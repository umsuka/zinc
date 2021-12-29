package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	zerolog "github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/routes"
	"github.com/prabhatsharma/zinc/pkg/zutils"

	"github.com/pyroscope-io/client/pyroscope"
)

func main() {

	serverAddress := os.Getenv("PYROSCOPE_SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "http://localhost:4040"
	}
	pyroscope.Start(pyroscope.Config{
		ApplicationName: "zinc.app",
		ServerAddress:   serverAddress,
		Logger:          pyroscope.StandardLogger,
	})

	err := godotenv.Load()
	if err != nil {
		zerolog.Print("Error loading .env file")
	}

	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	routes.SetRoutes(r) // Set up all API routes.

	// Run the server

	PORT := zutils.GetEnv("PORT", "4080")

	r.Run(":" + PORT)
}
