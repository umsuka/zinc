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
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/meta"
)

func TestIndex_Search(t *testing.T) {

	type args struct {
		iQuery *meta.ZincQuery
	}
	tests := []struct {
		name    string
		args    args
		data    []map[string]interface{}
		want    *meta.SearchResponse
		wantErr bool
	}{
		{
			name: "Search Query - Match",
			args: args{
				iQuery: &meta.ZincQuery{
					Query: &meta.Query{
						Match: map[string]*meta.MatchQuery{
							"_all": {
								Query: "Prabhat",
							},
						},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - Term",
			args: args{
				iQuery: &meta.ZincQuery{
					Query: &meta.Query{
						Term: map[string]*meta.TermQuery{
							"_all": {
								Value: "angeles",
							},
						},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - MatchAll",
			args: args{
				iQuery: &meta.ZincQuery{
					Query: &meta.Query{
						MatchAll: &meta.MatchAllQuery{},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - wildcard",
			args: args{
				iQuery: &meta.ZincQuery{
					Query: &meta.Query{
						Wildcard: map[string]*meta.WildcardQuery{
							"_all": {
								Value: "san*",
							},
						},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - fuzzy",
			args: args{
				iQuery: &meta.ZincQuery{
					Query: &meta.Query{
						Fuzzy: map[string]*meta.FuzzyQuery{
							"_all": {
								Value: "fransisco", // note the wrong spelling
							},
						},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - querystring1",
			args: args{
				iQuery: &meta.ZincQuery{
					Query: &meta.Query{
						QueryString: &meta.QueryStringQuery{
							Query: "angeles",
						},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - highlight",
			args: args{
				iQuery: &meta.ZincQuery{
					Query: &meta.Query{
						QueryString: &meta.QueryStringQuery{
							Query: "angeles",
						},
					},
					Timeout: 1,
					Fields:  []interface{}{"address.city"},
					Highlight: &meta.Highlight{
						Fields: map[string]*meta.Highlight{
							"address.city": {
								PreTags:  []string{"<b>"},
								PostTags: []string{"</b>"},
							},
						},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
	}

	indexName := "Search.index_1"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := NewIndex(indexName, "disk", nil)
			assert.Nil(t, err)
			assert.NotNil(t, index)

			if (index.CachedMappings) == nil {
				index.CachedMappings = meta.NewMappings()
			}
			index.CachedMappings.Properties["address.city"] = meta.Property{
				Type:          "text",
				Index:         true,
				Store:         true,
				Highlightable: true,
			}

			for _, d := range tt.data {
				rand.Seed(time.Now().UnixNano())
				docId := rand.Intn(1000)
				index.UpdateDocument(strconv.Itoa(docId), d, true)
			}
			got, err := index.Search(tt.args.iQuery)
			assert.Nil(t, err)
			assert.Equal(t, 1, got.Hits.Total.Value)

			DeleteIndex(indexName)
		})
	}

}
