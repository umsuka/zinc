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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildBlugeDocumentFromJSON(t *testing.T) {
	indexName := "TestBuildBlugeDocumentFromJSON.index_1"
	t.Run("build bluge document from json", func(t *testing.T) {
		idx, err := NewIndex(indexName, "disk", nil)
		assert.Nil(t, err)

		doc1 := make(map[string]interface{})
		doc1["id"] = "1"
		doc1["name"] = "test1"
		doc1["age"] = 10
		doc1["address"] = map[string]interface{}{
			"street": "447 Great Mall Dr",
			"city":   "Milpitas",
			"state":  "CA",
			"zip":    "95035",
		}

		doc, err := idx.BuildBlugeDocumentFromJSON("1", doc1)
		assert.Nil(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, "1", string(doc.ID().Term()))
	})

	t.Run("cleanup", func(t *testing.T) {
		err := DeleteIndex(indexName)
		assert.Nil(t, err)
	})
}
