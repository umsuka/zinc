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

package zinc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFrontendAssets(t *testing.T) {
	f, err := GetFrontendAssets()
	assert.Nil(t, err)
	assert.NotNil(t, f)
	t.Run("index.html", func(t *testing.T) {
		ff, err := f.Open("index.html")
		assert.Nil(t, err)
		fs, err := ff.Stat()
		assert.Nil(t, err)
		assert.Equal(t, "index.html", fs.Name())
	})
}
