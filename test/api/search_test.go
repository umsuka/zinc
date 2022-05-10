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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	meta "github.com/zinclabs/zinc/pkg/meta/v1"
)

func TestSearch(t *testing.T) {

	Convey("POST /api/:target/_search", t, func() {
		Convey("init data for search", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("search document with not exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{}`)
			resp := request("POST", "/api/notExistSearch/_search", body)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("search document with exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "alldocuments"}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("search document with not exist term", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "match", "query": {"term": "xxxx"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldEqual, 0)
		})
		Convey("search document with exist term", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: alldocuments", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "alldocuments", "query": {}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: wildcard", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "wildcard", "query": {"term": "dem*"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: fuzzy", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "fuzzy", "query": {"term": "demtschenk"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: term", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "term", 
				"query": {
					"term": "Turin", 
					"field":"City"
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: daterange", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(fmt.Sprintf(`{
				"search_type": "daterange",
				"query": {
					"start_time": "%s",
					"end_time": "%s"
				}
			}`,
				time.Now().UTC().Add(time.Hour*-24).Format("2006-01-02T15:04:05Z"),
				time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			))
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: matchall", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "matchall", "query": {"term": "demtschenk"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: match", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: matchphrase", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "matchphrase", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: multiphrase", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "multiphrase",
				"query": {
					"terms": [
						["demtschenko"],
						["albert"]
					]
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: prefix", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "prefix", "query": {"term": "dem"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: querystring", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "querystring", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
	})

	Convey("POST /api/:target/_search with aggregations", t, func() {
		Convey("terms aggregation", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "matchall", 
				"aggs": {
					"my-agg": {
						"agg_type": "terms",
						"field": "City"
					}
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(len(data.Aggregations), ShouldBeGreaterThanOrEqualTo, 1)
		})

		Convey("metric aggregation", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "matchall", 
				"aggs": {
					"my-agg-max": {
						"agg_type": "max",
						"field": "Year"
					},
					"my-agg-min": {
						"agg_type": "min",
						"field": "Year"
					},
					"my-agg-avg": {
						"agg_type": "avg",
						"field": "Year"
					}
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(len(data.Aggregations), ShouldBeGreaterThanOrEqualTo, 1)
		})
	})

}
