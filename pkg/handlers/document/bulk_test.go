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

package document

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/core"
)

func TestBulkWorker(t *testing.T) {
	input := `{ "index" : { "_index" : "olympics" } } 
	{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HAJOS, Alfred", "Country": "HUN", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Gold", "Season": "summer"}
	{ "index" : { "_index" : "olympics" } } 
	{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HERSCHMANN, Otto", "Country": "AUT", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Silver", "Season": "summer"}`

	rc := strings.NewReader("")

	// TODO add more test units:
	// 1. with delete
	// 2. check total num after bulk
	// 3. add error unit, without index, or without _id

	type args struct {
		target string
		body   *strings.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *BulkResponse
		wantErr bool
	}{
		{
			name: "bulk-with-index-name",
			args: args{
				target: "olympics",
				body:   rc,
			},
		},
		{
			name: "bulk-without-index-name",
			args: args{
				target: "",
				body:   rc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.body.Reset(input)
			got, err := BulkWorker(tt.args.target, tt.args.body)
			assert.Nil(t, err)
			assert.Equal(t, len(got.Items), 2)
			assert.Equal(t, got.Items[0]["index"].Status, 200)
			assert.Equal(t, got.Items[1]["index"].Status, 200)
		})
	}

	core.DeleteIndex("olympics")
}
