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
	"math/rand"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewIndex(t *testing.T) {
	Convey("test new index storage - disk", t, func() {
		rand.Seed(time.Now().UnixNano())
		id := rand.Intn(1000)
		indexName := "create.new.index_" + strconv.Itoa(id)

		index, err := NewIndex(indexName, "disk", UseNewIndexMeta, nil)
		So(err, ShouldBeNil)
		So(index.Name, ShouldEqual, indexName)
	})
	Convey("test new index storage s3", t, func() {
		// TODO: support
	})
	Convey("test new index storage minio", t, func() {
		// TODO: support
	})

	// Cleanup data folder
	// err := os.RemoveAll("data")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
