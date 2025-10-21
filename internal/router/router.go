package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Hapaa16/go-shortener/internal/handlers"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// health
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	uh := handlers.NewUrlHandler(db)
	api := r.Group("/api")
	{
		api.GET("/shorten/:url", uh.AccessUrl)
		api.GET("/shorten/:url/stats", uh.GetStats)
		api.POST("/shorten", uh.ShortenUrl)
		api.PUT("/shorten/:url", uh.UpdateUrl)
		api.DELETE("/shorten/:url", uh.DeleteUrl)

	}

	return r
}
