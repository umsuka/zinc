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

package search

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/errors"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
)

// SearchDSL searches the index for the given http request from end user
func SearchDSL(c *gin.Context) {
	indexName := c.Param("target")

	query := new(meta.ZincQuery)
	if err := c.BindJSON(query); err != nil {
		log.Printf("handlers.search.searchDSL: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storageType := "disk"
	indexSize := 0.0

	var err error
	var resp *meta.SearchResponse
	if indexName == "" || strings.HasSuffix(indexName, "*") {
		resp, err = core.MultiSearchV2(indexName, query)
	} else {
		index, exists := core.GetIndex(indexName)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
			return
		}

		storageType = index.StorageType
		indexSize = index.StorageSize
		resp, err = index.SearchV2(query)
	}

	if err != nil {
		errors.HandleError(c, err)
		return
	}

	eventData := make(map[string]interface{})
	eventData["search_type"] = "query_dsl"
	eventData["search_index_storage"] = storageType
	eventData["search_index_size_in_mb"] = indexSize
	eventData["time_taken_to_search_in_ms"] = resp.Took
	eventData["aggregations_count"] = len(query.Aggregations)
	core.Telemetry.Event("search", eventData)

	c.JSON(http.StatusOK, resp)
}
