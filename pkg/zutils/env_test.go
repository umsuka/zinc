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

	. "github.com/smartystreets/goconvey/convey"
)

func TestEnv(t *testing.T) {
	Convey("zutils:env", t, func() {
		Convey("GetEnv", func() {
			a := GetEnv("ZINC_SENTRY", "true")
			So(a, ShouldEqual, "true")
			a = GetEnv("ZINC_SENTRY", "")
			So(a, ShouldEqual, "")
		})
		Convey("GetEnvToUpper", func() {
			a := GetEnvToUpper("ZINC_SENTRY", "true")
			So(a, ShouldEqual, "TRUE")
		})
		Convey("GetEnvToLower", func() {
			a := GetEnvToLower("ZINC_SENTRY", "TRUE")
			So(a, ShouldEqual, "true")
		})
		Convey("GetEnvBool", func() {
			a := GetEnvToBool("ZINC_SENTRY", "true")
			So(a, ShouldEqual, true)
			a = GetEnvToBool("ZINC_SENTRY", "")
			So(a, ShouldEqual, false)
		})
	})
}
