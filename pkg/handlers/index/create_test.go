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

package index

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/zinclabs/zinc/pkg/core"
)

func TestCreateIndexWorker(t *testing.T) {
	type args struct {
		newIndex  *core.Index
		indexName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		// 1. create index simple
		// 2. create index with mappings
		// 3. create index with settings
		// 4. create index with mappings and settings
		// 5. create index use template
		// 6. create index without index name
		// 7. create index with error settings
		{
			name: "test1",
			args: args{
				newIndex: &core.Index{
					StorageType: "disk",
				},
				indexName: "test1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(time.Now().UnixNano())
			id := rand.Intn(1000)

			if err := CreateIndexWorker(tt.args.newIndex, tt.args.indexName+strconv.Itoa(id)); (err != nil) != tt.wantErr {
				t.Errorf("CreateIndexWorker() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := core.DeleteIndex(tt.args.indexName + strconv.Itoa(id)); err != nil {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
