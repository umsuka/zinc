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

package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestApiAuth(t *testing.T) {
	Convey("test auth api", t, func() {
		r := server()
		Convey("check auth with auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, password)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("check auth with error password", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, "xxx")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("check auth without auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})
}
