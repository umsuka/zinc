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

package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
	"github.com/zinclabs/zinc/pkg/uquery/v2/mappings"
)

func GetIndexMapping(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	// format mappings
	mappings := index.CachedMappings
	if mappings == nil {
		mappings = meta.NewMappings()
	}

	c.JSON(http.StatusOK, gin.H{index.Name: gin.H{"mappings": mappings}})
}

func UpdateIndexMapping(c *gin.Context) {
	indexName := c.Param("target")
	if indexName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index.name should be not empty"})
		return
	}

	var newIndex core.Index
	if err := c.BindJSON(&newIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var index *core.Index
	index, exists := core.GetIndex(indexName)
	if !exists {
		index1, err := core.NewIndex(indexName, newIndex.StorageType, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		index = index1
	}

	mappings, err := mappings.Request(nil, newIndex.Mappings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update mappings
	if mappings != nil && len(mappings.Properties) > 0 {
		index.SetMappings(mappings)
	}

	// store index
	core.StoreIndex(index)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
