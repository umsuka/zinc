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

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

// Search searches the index for the given http request from end user
func Search(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	var iQuery v1.ZincQuery
	if err := c.BindJSON(&iQuery); err != nil {
		log.Printf("handlers.search.Search: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := index.Search(&iQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventData := make(map[string]interface{})
	eventData["search_type"] = iQuery.SearchType
	eventData["search_index_storage"] = index.StorageType
	eventData["search_index_size_in_mb"] = index.StorageSize
	eventData["time_taken_to_search_in_ms"] = res.Took
	eventData["aggregations_count"] = len(iQuery.Aggregations)
	core.Telemetry.Event("search", eventData)

	c.JSON(http.StatusOK, res)
}
