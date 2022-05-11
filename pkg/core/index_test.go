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
	"reflect"
	"testing"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"
	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/meta"
)

func TestIndex_BuildBlugeDocumentFromJSON(t *testing.T) {
	var index *Index
	var err error
	indexName := "TestIndex_BuildBlugeDocumentFromJSON.index_1"

	type args struct {
		docID string
		doc   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		init    func()
		want    *bluge.Document
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				docID: "1",
				doc: map[string]interface{}{
					"id":     "1",
					"name":   "test1",
					"age":    10,
					"length": 3.14,
					"dev":    true,
					"address": map[string]interface{}{
						"street": "447 Great Mall Dr",
						"city":   "Milpitas",
						"state":  "CA",
						"zip":    95035,
					},
					"tag1":       []interface{}{"tag1", "tag2"},
					"tag2":       []interface{}{3.14, 3.15},
					"tag3":       []interface{}{true, false},
					"@timestamp": time.Now().Format(time.RFC3339),
					"time":       time.Now().Format(time.RFC3339),
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: false,
		},
		{
			name: "timestamp with epoch_millis",
			args: args{
				docID: "2",
				doc: map[string]interface{}{
					"id":         "2",
					"name":       "test1",
					"age":        10,
					"length":     3.14,
					"dev":        true,
					"@timestamp": float64(1652176732575),
					"time":       float64(1652176732575),
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: false,
		},
		{
			name: "with analyzer",
			args: args{
				docID: "3",
				doc: map[string]interface{}{
					"id":     "3",
					"name":   3,
					"age":    "10",
					"length": 3,
					"dev":    true,
				},
			},
			init: func() {
				index.CachedMappings.Properties["id"] = meta.Property{
					Type:          "keyword",
					Index:         true,
					Store:         true,
					Highlightable: true,
				}
				index.CachedMappings.Properties["name"] = meta.Property{
					Type:     "text",
					Index:    true,
					Analyzer: "default",
				}
				index.CachedAnalyzers["default"] = analyzer.NewStandardAnalyzer()
			},
			want:    &bluge.Document{},
			wantErr: true,
		},
		{
			name: "type conflict text",
			args: args{
				docID: "4",
				doc: map[string]interface{}{
					"id":   "4",
					"name": 3,
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: true,
		},
		{
			name: "type conflict numeric",
			args: args{
				docID: "4",
				doc: map[string]interface{}{
					"id":     "4",
					"name":   "test1",
					"age":    "10",
					"length": 3,
					"dev":    true,
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: true,
		},
		{
			name: "keyword type float64",
			args: args{
				docID: "4",
				doc: map[string]interface{}{
					"id": 3.14,
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: false,
		},
		{
			name: "keyword type int",
			args: args{
				docID: "4",
				doc: map[string]interface{}{
					"id": 3,
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: false,
		},
		{
			name: "keyword type bool",
			args: args{
				docID: "4",
				doc: map[string]interface{}{
					"id": false,
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: false,
		},
		{
			name: "keyword type other",
			args: args{
				docID: "4",
				doc: map[string]interface{}{
					"id": []byte("foo"),
				},
			},
			init:    func() {},
			want:    &bluge.Document{},
			wantErr: false,
		},
	}

	t.Run("prepare", func(t *testing.T) {
		index, err = NewIndex(indexName, "disk", nil)
		assert.Nil(t, err)
		assert.NotNil(t, index)
		StoreIndex(index)
		index.CachedMappings.Properties["time"] = meta.NewProperty("time")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init()
			got, err := index.BuildBlugeDocumentFromJSON(tt.args.docID, tt.args.doc)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NotNil(t, got)
			wantType := reflect.TypeOf(tt.want)
			gotType := reflect.TypeOf(got)
			assert.Equal(t, wantType, gotType)
		})
	}

	t.Run("cleanup", func(t *testing.T) {
		err := DeleteIndex(indexName)
		assert.Nil(t, err)
	})
}

func TestIndex_Settings(t *testing.T) {
	type fields struct {
		Name                string
		IndexType           string
		StorageType         string
		DocsCount           int64
		StorageSize         float64
		StorageSizeNextTime time.Time
		Mappings            map[string]interface{}
		Settings            *meta.IndexSettings
		CachedAnalyzers     map[string]*analysis.Analyzer
		CachedMappings      *meta.Mappings
		Writer              *bluge.Writer
	}
	type args struct {
		settings *meta.IndexSettings
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := &Index{
				Name:                tt.fields.Name,
				IndexType:           tt.fields.IndexType,
				StorageType:         tt.fields.StorageType,
				DocsCount:           tt.fields.DocsCount,
				StorageSize:         tt.fields.StorageSize,
				StorageSizeNextTime: tt.fields.StorageSizeNextTime,
				Mappings:            tt.fields.Mappings,
				Settings:            tt.fields.Settings,
				CachedAnalyzers:     tt.fields.CachedAnalyzers,
				CachedMappings:      tt.fields.CachedMappings,
				Writer:              tt.fields.Writer,
			}
			if err := index.SetSettings(tt.args.settings); (err != nil) != tt.wantErr {
				t.Errorf("Index.SetSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
