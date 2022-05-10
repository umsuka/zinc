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

func TestCreateUser(t *testing.T) {
	type args struct {
		userID            string
		name              string
		plaintextPassword string
		role              string
	}
	tests := []struct {
		name    string
		args    args
		want    *ZincUser
		wantErr bool
	}{
		{
			name: "create user",
			args: args{
				userID:            "testuser",
				name:              "Test User",
				plaintextPassword: "testpassword",
				role:              "admin",
			},
			want: &ZincUser{
				ID:   "testuser",
				Name: "Test User",
				Role: "admin",
			},
			wantErr: false,
		},
		{
			name: "update exists user",
			args: args{
				userID:            "testuser",
				name:              "Test User Updated",
				plaintextPassword: "testpassword",
				role:              "admin",
			},
			want: &ZincUser{
				ID:   "testuser",
				Name: "Test User Updated",
				Role: "admin",
			},
			wantErr: false,
		},
		{
			name: "create user with empty userID",
			args: args{
				userID: "",
			},
			want: &ZincUser{
				ID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUser(tt.args.userID, tt.args.name, tt.args.plaintextPassword, tt.args.role)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)

			salt := got.Salt
			password := GeneratePassword(tt.args.plaintextPassword, salt)
			assert.Equal(t, password, got.Password)
		})
	}
}
