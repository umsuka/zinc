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
	"github.com/zinclabs/zinc/pkg/handlers/auth"
	"github.com/zinclabs/zinc/pkg/handlers/document"
	"github.com/zinclabs/zinc/pkg/handlers/index"
	"github.com/zinclabs/zinc/pkg/handlers/search"
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

	// UI
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

	// auth
	r.POST("/api/login", auth.Login)
	r.PUT("/api/user", AuthMiddleware, auth.CreateUpdate)
	r.DELETE("/api/user/:id", AuthMiddleware, auth.Delete)
	r.GET("/api/user", AuthMiddleware, auth.List)

	// index
	r.GET("/api/index", AuthMiddleware, index.List)
	r.PUT("/api/index", AuthMiddleware, index.Create)
	r.DELETE("/api/index/:target", AuthMiddleware, index.Delete)
	// index settings
	r.GET("/api/:target/_mapping", AuthMiddleware, index.GetMapping)
	r.PUT("/api/:target/_mapping", AuthMiddleware, index.SetMapping)
	r.GET("/api/:target/_settings", AuthMiddleware, index.GetSettings)
	r.PUT("/api/:target/_settings", AuthMiddleware, index.SetSettings)
	r.POST("/api/:target/_analyze", AuthMiddleware, index.Analyze)
	r.POST("/api/_analyze", AuthMiddleware, index.Analyze)

	// document

	// Document Bulk update/insert
	r.POST("/api/_bulk", AuthMiddleware, document.Bulk)
	// Document CRUD APIs. Update is same as create.
	r.PUT("/api/:target/_doc", AuthMiddleware, document.CreateUpdate)
	r.PUT("/api/:target/_doc/:id", AuthMiddleware, document.CreateUpdate)
	r.DELETE("/api/:target/_doc/:id", AuthMiddleware, document.Delete)

	// search
	r.POST("/api/:target/_search", AuthMiddleware, search.Search)

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

	r.POST("/es/_search", AuthMiddleware, search.SearchDSL)
	r.POST("/es/:target/_search", AuthMiddleware, search.SearchDSL)

	r.GET("/es/_index_template", AuthMiddleware, index.ListTemplate)
	r.PUT("/es/_index_template/:target", AuthMiddleware, index.UpdateTemplate)
	r.GET("/es/_index_template/:target", AuthMiddleware, index.GetTemplate)
	r.HEAD("/es/_index_template/:target", AuthMiddleware, index.GetTemplate)
	r.DELETE("/es/_index_template/:target", AuthMiddleware, index.DeleteTemplate)

	r.GET("/es/:target/_mapping", AuthMiddleware, index.GetMapping)
	r.PUT("/es/:target/_mapping", AuthMiddleware, index.SetMapping)

	r.GET("/es/:target/_settings", AuthMiddleware, index.GetSettings)
	r.PUT("/es/:target/_settings", AuthMiddleware, index.SetSettings)

	r.POST("/es/_analyze", AuthMiddleware, index.Analyze)
	r.POST("/es/:target/_analyze", AuthMiddleware, index.Analyze)

	// ES Bulk update/insert
	r.POST("/es/_bulk", AuthMiddleware, document.ESBulk)
	r.POST("/es/:target/_bulk", AuthMiddleware, document.ESBulk)
	// ES Document
	r.POST("/es/:target/_doc", AuthMiddleware, document.CreateUpdate)
	r.PUT("/es/:target/_doc/:id", AuthMiddleware, document.CreateUpdate)
	r.PUT("/es/:target/_create/:id", AuthMiddleware, document.CreateUpdate)
	r.POST("/es/:target/_create/:id", AuthMiddleware, document.CreateUpdate)
	r.POST("/es/:target/_update/:id", AuthMiddleware, document.CreateUpdate)
	r.DELETE("/es/:target/_doc/:id", AuthMiddleware, document.Delete)

	core.Telemetry.Instance()
	core.Telemetry.Event("server_start", nil)
	core.Telemetry.Cron()
}
