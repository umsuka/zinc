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
	"net/http"

	"github.com/gin-gonic/gin"
)

func ZincAuthMiddleware(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth {
		if _, ok := VerifyCredentials(user, password); ok {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth": "Invalid credentials"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth": "Missing credentials"})
		return
	}
}

func VerifyCredentials(userID, password string) (SimpleUser, bool) {
	user, ok := ZINC_CACHED_USERS[userID]
	if !ok {
		return SimpleUser{}, false
	}

	incomingEncryptedPassword := GeneratePassword(password, user.Salt)
	if incomingEncryptedPassword == user.Password {
		return user, true
	}

	return SimpleUser{}, false
}
