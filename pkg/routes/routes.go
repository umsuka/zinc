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

package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc"
	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/handlers"
	handlersV2 "github.com/zinclabs/zinc/pkg/handlers/v2"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

// SetRoutes sets up all gin HTTP API endpoints that can be called by front end
func SetRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// debug for accesslog
	if gin.Mode() == gin.DebugMode {
		AccessLog(r)
	}

	r.GET("/", v1.GUI)
	r.GET("/version", v1.GetVersion)
	r.GET("/healthz", v1.GetHealthz)

	front, err := zinc.GetFrontendAssets()
	if err != nil {
		log.Err(err)
	}

	HTTPCacheForUI(r)
	r.StaticFS("/ui/", http.FS(front))
	r.NoRoute(func(c *gin.Context) {
		log.Error().
			Str("method", c.Request.Method).
			Int("code", 404).
			Int("took", 0).
			Msg(c.Request.RequestURI)

		if strings.HasPrefix(c.Request.RequestURI, "/ui/") {
			path := strings.TrimPrefix(c.Request.RequestURI, "/ui/")
			locationPath := strings.Repeat("../", strings.Count(path, "/"))
			c.Status(http.StatusFound)
			c.Writer.Header().Set("Location", "./"+locationPath)
		}
	})

	r.POST("/api/login", handlers.ValidateCredentials)

	r.PUT("/api/user", AuthMiddleware, handlers.CreateUpdateUser)
	r.DELETE("/api/user/:userID", AuthMiddleware, handlers.DeleteUser)
	r.GET("/api/users", AuthMiddleware, handlers.GetUsers)

	r.GET("/api/index", AuthMiddleware, handlers.ListIndexes)
	r.PUT("/api/index", AuthMiddleware, handlers.CreateIndex)
	r.PUT("/api/index/:target", AuthMiddleware, handlers.CreateIndex)
	r.DELETE("/api/index/:target", AuthMiddleware, handlers.DeleteIndex)

	// Bulk update/insert
	r.POST("/api/_bulk", AuthMiddleware, handlers.BulkHandler)
	r.POST("/api/:target/_bulk", AuthMiddleware, handlers.BulkHandler)

	// Document CRUD APIs. Update is same as create.
	r.PUT("/api/:target/document", AuthMiddleware, handlers.UpdateDocument)
	r.POST("/api/:target/_doc", AuthMiddleware, handlers.UpdateDocument)
	r.PUT("/api/:target/_doc/:id", AuthMiddleware, handlers.UpdateDocument)
	r.POST("/api/:target/_search", AuthMiddleware, handlers.SearchIndex)
	r.DELETE("/api/:target/_doc/:id", AuthMiddleware, handlers.DeleteDocument)

	r.GET("/api/:target/_mapping", AuthMiddleware, handlersV2.GetIndexMapping)
	r.PUT("/api/:target/_mapping", AuthMiddleware, handlersV2.UpdateIndexMapping)

	r.GET("/api/:target/_settings", AuthMiddleware, handlersV2.GetIndexSettings)
	r.PUT("/api/:target/_settings", AuthMiddleware, handlersV2.UpdateIndexSettings)

	r.POST("/api/_analyze", AuthMiddleware, handlersV2.Analyze)
	r.POST("/api/:target/_analyze", AuthMiddleware, handlersV2.Analyze)

	/**
	 * elastic compatible APIs
	 */

	r.GET("/es/", func(c *gin.Context) {
		c.JSON(http.StatusOK, v1.NewESInfo(c))
	})
	r.GET("/es/_license", func(c *gin.Context) {
		c.JSON(http.StatusOK, v1.NewESLicense(c))
	})
	r.GET("/es/_xpack", func(c *gin.Context) {
		c.JSON(http.StatusOK, v1.NewESXPack(c))
	})

	r.POST("/es/_search", AuthMiddleware, handlersV2.SearchIndex)
	r.POST("/es/:target/_search", AuthMiddleware, handlersV2.SearchIndex)

	r.GET("/es/_index_template", AuthMiddleware, handlersV2.ListIndexTemplate)
	r.PUT("/es/_index_template/:target", AuthMiddleware, handlersV2.UpdateIndexTemplate)
	r.GET("/es/_index_template/:target", AuthMiddleware, handlersV2.GetIndexTemplate)
	r.HEAD("/es/_index_template/:target", AuthMiddleware, handlersV2.GetIndexTemplate)
	r.DELETE("/es/_index_template/:target", AuthMiddleware, handlersV2.DeleteIndexTemplate)

	r.GET("/es/:target/_mapping", AuthMiddleware, handlersV2.GetIndexMapping)
	r.PUT("/es/:target/_mapping", AuthMiddleware, handlersV2.UpdateIndexMapping)

	r.GET("/es/:target/_settings", AuthMiddleware, handlersV2.GetIndexSettings)
	r.PUT("/es/:target/_settings", AuthMiddleware, handlersV2.UpdateIndexSettings)

	r.POST("/es/_analyze", AuthMiddleware, handlersV2.Analyze)
	r.POST("/es/:target/_analyze", AuthMiddleware, handlersV2.Analyze)

	r.POST("/es/:target/_doc", AuthMiddleware, handlers.UpdateDocument)
	r.PUT("/es/:target/_doc/:id", AuthMiddleware, handlers.UpdateDocument)
	r.PUT("/es/:target/_create/:id", AuthMiddleware, handlers.UpdateDocument)
	r.POST("/es/:target/_create/:id", AuthMiddleware, handlers.UpdateDocument)
	r.POST("/es/:target/_update/:id", AuthMiddleware, handlers.UpdateDocument)
	r.DELETE("/es/:target/_doc/:id", AuthMiddleware, handlers.DeleteDocument)

	// Bulk update/insert
	r.POST("/es/_bulk", AuthMiddleware, handlers.ESBulkHandler)
	r.POST("/es/:target/_bulk", AuthMiddleware, handlers.ESBulkHandler)

	core.Telemetry.Instance()
	core.Telemetry.Event("server_start", nil)
	core.Telemetry.Cron()
}
