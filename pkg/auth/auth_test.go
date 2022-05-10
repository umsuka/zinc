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

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyCredentials(t *testing.T) {
	type args struct {
		userID   string
		password string
	}
	tests := []struct {
		name  string
		args  args
		want  SimpleUser
		want1 bool
	}{
		{
			name: "test with admin",
			args: args{
				userID:   "admin",
				password: "Complexpass#123",
			},
			want: SimpleUser{
				ID: "admin",
			},
			want1: true,
		},
		{
			name: "test with error password",
			args: args{
				userID:   "admin",
				password: "xxxxxxxx",
			},
			want: SimpleUser{
				ID: "",
			},
			want1: false,
		},
		{
			name: "test with error user",
			args: args{
				userID:   "xxxxxxxx",
				password: "xxxxxxxx",
			},
			want: SimpleUser{
				ID: "",
			},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := VerifyCredentials(tt.args.userID, tt.args.password)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
