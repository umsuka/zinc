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

func TestLoadIndexes(t *testing.T) {
	t.Run("load system index", func(t *testing.T) {
		// index cann't be reopen, so need close first
		for _, index := range ZINC_SYSTEM_INDEX_LIST {
			index.Writer.Close()
		}
		var err error
		ZINC_SYSTEM_INDEX_LIST, err = LoadZincSystemIndexes()
		assert.Nil(t, err)
		assert.Equal(t, len(systemIndexList), len(ZINC_SYSTEM_INDEX_LIST))
		assert.Equal(t, "_index_mapping", ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Name)
	})

	t.Run("load user index from disk", func(t *testing.T) {
		// index cann't be reopen, so need close first
		for _, index := range ZINC_INDEX_LIST {
			index.Writer.Close()
		}
		var err error
		ZINC_INDEX_LIST, err = LoadZincIndexesFromMeta()
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, 0, len(ZINC_INDEX_LIST))
	})
}
