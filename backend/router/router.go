package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"autopowerhub/api/handler"
	"autopowerhub/middleware"
	authsvc "autopowerhub/service/auth"

	"github.com/gin-gonic/gin"
)

func Setup(
	authH *handler.AuthHandler,
	deviceH *handler.DeviceHandler,
	debugH *handler.DebugHandler,
	authSvc *authsvc.Service,
	baseDir string, // directory containing config.yaml; all static paths are relative to this
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/login", authH.Login)

		protected := api.Group("")
		protected.Use(middleware.JWT(authSvc))
		{
			protected.GET("/device", deviceH.ListDevices)
			protected.POST("/device/:id/power", deviceH.Power)
			protected.POST("/device/:id/test", deviceH.Test)
			protected.GET("/debug/device/:id/scan", debugH.BLEScan)
		}
	}

	// Serve Vue SPA from <baseDir>/web/ (produced by `cd frontend && npm run build`)
	webDir := filepath.Join(baseDir, "web")
	indexHTML := filepath.Join(webDir, "index.html")

	r.Static("/assets", filepath.Join(webDir, "assets"))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
			return
		}
		if _, err := os.Stat(indexHTML); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "frontend not built yet",
				"hint":  "cd frontend && npm install && npm run build",
			})
			return
		}
		c.File(indexHTML)
	})

	return r
}
