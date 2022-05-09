/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package core

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	meta "github.com/zinclabs/zinc/pkg/meta/v2"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/v2/analysis"
	"github.com/zinclabs/zinc/pkg/zutils"
	"github.com/zinclabs/zinc/pkg/zutils/flatten"
)

// BuildBlugeDocumentFromJSON returns the bluge document for the json document. It also updates the mapping for the fields if not found.
// If no mappings are found, it creates te mapping for all the encountered fields. If mapping for some fields is found but not for others
// then it creates the mapping for the missing fields.
func (index *Index) BuildBlugeDocumentFromJSON(docID string, doc map[string]interface{}) (*bluge.Document, error) {
	// Pick the index mapping from the cache if it already exists
	mappings := index.CachedMappings
	if mappings == nil {
		mappings = meta.NewMappings()
	}

	mappingsNeedsUpdate := false

	// Create a new bluge document
	bdoc := bluge.NewDocument(docID)
	flatDoc, _ := flatten.Flatten(doc, "")
	// Iterate through each field and add it to the bluge document
	for key, value := range flatDoc {
		if value == nil || key == "@timestamp" {
			continue
		}

		if _, ok := mappings.Properties[key]; !ok {
			// try to find the type of the value and use it to define default mapping
			switch value.(type) {
			case string:
				mappings.Properties[key] = meta.NewProperty("text")
			case float64:
				mappings.Properties[key] = meta.NewProperty("numeric")
			case bool:
				mappings.Properties[key] = meta.NewProperty("bool")
			case []interface{}:
				if v, ok := value.([]interface{}); ok {
					for _, vv := range v {
						switch vv.(type) {
						case string:
							mappings.Properties[key] = meta.NewProperty("text")
						case float64:
							mappings.Properties[key] = meta.NewProperty("numeric")
						case bool:
							mappings.Properties[key] = meta.NewProperty("bool")
						}
						break
					}
				}
			}

			mappingsNeedsUpdate = true
		}

		if !mappings.Properties[key].Index {
			continue // not index, skip
		}

		switch v := value.(type) {
		case []interface{}:
			for _, v := range v {
				if err := index.buildField(mappings, bdoc, key, v); err != nil {
					return nil, err
				}
			}
		default:
			if err := index.buildField(mappings, bdoc, key, v); err != nil {
				return nil, err
			}
		}
	}

	if mappingsNeedsUpdate {
		index.SetMappings(mappings)
		StoreIndex(index)
	}

	timestamp := time.Now()
	if v, ok := flatDoc["@timestamp"]; ok {
		switch v := v.(type) {
		case string:
			if t, err := time.Parse(time.RFC3339, v); err == nil && !t.IsZero() {
				timestamp = t
				delete(doc, "@timestamp")
			}
		case float64:
			if t := zutils.Unix(int64(v)); !t.IsZero() {
				timestamp = t
				delete(doc, "@timestamp")
			}
		default:
			// noop
		}
	}
	docByteVal, _ := json.Marshal(doc)
	bdoc.AddField(bluge.NewDateTimeField("@timestamp", timestamp).StoreValue().Sortable().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_index", []byte(index.Name)))
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", []string{"_index", "_id", "_source", "@timestamp"}))

	return bdoc, nil
}

func (index *Index) buildField(mappings *meta.Mappings, bdoc *bluge.Document, key string, value interface{}) error {
	var field *bluge.TermField
	switch mappings.Properties[key].Type {
	case "text":
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("field [%s] was set type to [text] but got a %T value", key, value)
		}
		field = bluge.NewTextField(key, v).SearchTermPositions()
		fieldAnalyzer, _ := zincanalysis.QueryAnalyzerForField(index.CachedAnalyzers, index.CachedMappings, key)
		if fieldAnalyzer != nil {
			field.WithAnalyzer(fieldAnalyzer)
		}
	case "numeric":
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("field [%s] was set type to [numeric] but got a %T value", key, value)
		}
		field = bluge.NewNumericField(key, v)
	case "keyword":
		// compatible verion <= v0.1.4
		if v, ok := value.(bool); ok {
			field = bluge.NewKeywordField(key, strconv.FormatBool(v))
		} else if v, ok := value.(string); ok {
			field = bluge.NewKeywordField(key, v)
		} else {
			return fmt.Errorf("keyword type only support text")
		}
	case "bool": // found using existing index mapping
		value := value.(bool)
		field = bluge.NewKeywordField(key, strconv.FormatBool(value))
	case "time":
		format := time.RFC3339
		if mappings.Properties[key].Format != "" {
			format = mappings.Properties[key].Format
		}
		tim, err := time.Parse(format, value.(string))
		if err != nil {
			return err
		}
		field = bluge.NewDateTimeField(key, tim)
	}

	if mappings.Properties[key].Store {
		field.StoreValue()
	}
	if mappings.Properties[key].Sortable {
		field.Sortable()
	}
	if mappings.Properties[key].Aggregatable {
		field.Aggregatable()
	}
	if mappings.Properties[key].Highlightable {
		field.HighlightMatches()
	}
	bdoc.AddField(field)

	return nil
}

func (index *Index) UseTemplate() error {
	template, err := UseTemplate(index.Name)
	if err != nil {
		return err
	}

	if template == nil {
		return nil
	}

	if template.Template.Settings != nil {
		index.SetSettings(template.Template.Settings)
	}

	if template.Template.Mappings != nil {
		index.SetMappings(template.Template.Mappings)
	}

	return nil
}

func (index *Index) SetSettings(settings *meta.IndexSettings) error {
	if settings == nil {
		return nil
	}

	index.Settings = settings

	return nil
}

func (index *Index) SetAnalyzers(analyzers map[string]*analysis.Analyzer) error {
	if len(analyzers) == 0 {
		return nil
	}

	index.CachedAnalyzers = analyzers

	return nil
}

func (index *Index) SetMappings(mappings *meta.Mappings) error {
	if mappings == nil || len(mappings.Properties) == 0 {
		return nil
	}

	// custom analyzer just for text field
	for _, prop := range mappings.Properties {
		if prop.Type != "text" {
			prop.Analyzer = ""
			prop.SearchAnalyzer = ""
		}
	}

	mappings.Properties["_id"] = meta.NewProperty("keyword")

	// @timestamp need date_range/date_histogram aggregation, and mappings used for type check in aggregation
	mappings.Properties["@timestamp"] = meta.NewProperty("time")

	// update in the cache
	index.CachedMappings = mappings
	index.Mappings = nil

	return nil
}

// DEPRECATED GetStoredMapping returns the mappings of all the indexes from _index_mapping system index
func (index *Index) GetStoredMapping() (*meta.Mappings, error) {
	log.Error().Bool("deprecated", true).Msg("GetStoredMapping is deprecated, use index.CachedMappings instead")
	for _, indexName := range systemIndexList {
		if index.Name == indexName {
			return nil, nil
		}
	}

	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Writer.Reader()
	defer reader.Close()

	// search for the index mapping _index_mapping index
	query := bluge.NewTermQuery(index.Name).SetField("_id")
	searchRequest := bluge.NewTopNSearch(1, query) // Should get just 1 result at max
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Error().Str("index", index.Name).Msg("error executing search: " + err.Error())
		return nil, err
	}

	next, err := dmi.Next()
	if err != nil {
		return nil, err
	}

	mappings := new(meta.Mappings)
	oldMappings := make(map[string]string)
	if next != nil {
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_source":
				if string(value) != "" {
					json.Unmarshal(value, mappings)
				}
			default:
				oldMappings[field] = string(value)
			}
			return true
		})
		if err != nil {
			return nil, err
		}
	}

	// compatible old mappings format
	if len(mappings.Properties) == 0 && len(oldMappings) > 0 {
		mappings.Properties = make(map[string]meta.Property, len(oldMappings))
		for k, v := range oldMappings {
			mappings.Properties[k] = meta.NewProperty(v)
		}
	}

	if len(mappings.Properties) == 0 {
		mappings.Properties = make(map[string]meta.Property)
	}

	return mappings, nil
}

func (index *Index) LoadDocsCount() (int64, error) {
	query := bluge.NewMatchAllQuery()
	searchRequest := bluge.NewTopNSearch(0, query).WithStandardAggregations()
	reader, _ := index.Writer.Reader()
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return 0, fmt.Errorf("core.index.LoadDocsCount: error executing search: %s", err.Error())
	}

	return int64(dmi.Aggregations().Count()), nil
}

func (index *Index) LoadStorageSize() float64 {
	size := 0.0
	switch index.StorageType {
	case "s3":
		return size // TODO: implement later
	case "minio":
		return size // TODO: implement later
	default:
		path := zutils.GetEnv("ZINC_DATA_PATH", "./data")
		indexLocation := filepath.Join(path, index.Name)
		size, _ = zutils.DirSize(indexLocation)
		return math.Round(size)
	}
}

func (index *Index) ReLoadStorageSize() {
	if index.StorageSizeNextTime.After(time.Now()) {
		return // skip
	}

	index.StorageSizeNextTime = time.Now().Add(time.Minute * 10)
	go func() {
		index.StorageSize = index.LoadStorageSize()
	}()
}

func (index *Index) ReduceDocsCount(n int64) {
	atomic.AddInt64(&index.DocsCount, -n)
	index.ReLoadStorageSize()
}

func (index *Index) GainDocsCount(n int64) {
	atomic.AddInt64(&index.DocsCount, n)
	index.ReLoadStorageSize()
}

func (index *Index) Close() error {
	return index.Writer.Close()
}
