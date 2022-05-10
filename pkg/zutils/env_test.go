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

package zutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	t.Run("GetEnv", func(t *testing.T) {
		a := GetEnv("ZINC_SENTRY", "true")
		assert.Equal(t, "true", a)
		a = GetEnv("ZINC_SENTRY", "")
		assert.Equal(t, "", a)
	})
	t.Run("GetEnvToUpper", func(t *testing.T) {
		a := GetEnvToUpper("ZINC_SENTRY", "true")
		assert.Equal(t, "TRUE", a)
	})
	t.Run("GetEnvToLower", func(t *testing.T) {
		a := GetEnvToLower("ZINC_SENTRY", "TRUE")
		assert.Equal(t, "true", a)
	})
	t.Run("GetEnvBool", func(t *testing.T) {
		a := GetEnvToBool("ZINC_SENTRY", "true")
		assert.Equal(t, true, a)
		a = GetEnvToBool("ZINC_SENTRY", "")
		assert.Equal(t, false, a)
	})
}
