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

package meta

import "time"

// SearchResponse for a query
type SearchResponse struct {
	Took         int                            `json:"took"` // Time it took to generate the response
	TimedOut     bool                           `json:"timed_out"`
	Shards       Shards                         `json:"_shards"`
	Hits         Hits                           `json:"hits"`
	Aggregations map[string]AggregationResponse `json:"aggregations,omitempty"`
	Error        string                         `json:"error"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total    Total   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Hit struct {
	Index     string                 `json:"_index"`
	Type      string                 `json:"_type"`
	ID        string                 `json:"_id"`
	Score     float64                `json:"_score"`
	Timestamp time.Time              `json:"@timestamp"`
	Source    interface{}            `json:"_source,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Highlight map[string]interface{} `json:"highlight,omitempty"`
}

type Total struct {
	Value int `json:"value"` // Count of documents returned
}

type AggregationResponse struct {
	Value    interface{} `json:"value,omitempty"`
	Buckets  interface{} `json:"buckets,omitempty"`  // slice or map
	Interval string      `json:"interval,omitempty"` // support for auto_date_histogram_aggregation
}

type AggregationBucket struct {
	Key          interface{}                    `json:"key"`
	KeyAsString  string                         `json:"key_as_string,omitempty"`
	DocCount     uint64                         `json:"doc_count"`
	Aggregations map[string]AggregationResponse `json:"aggregations,omitempty"`
}
