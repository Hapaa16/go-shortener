package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Hapaa16/go-shortener/internal/config"
	"github.com/Hapaa16/go-shortener/internal/db"
	"github.com/Hapaa16/go-shortener/internal/models"
	"github.com/Hapaa16/go-shortener/internal/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	d, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	autoMigrate(d)

	r := router.Setup(d)

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Println("server exiting")

}

func autoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Url{}); err != nil {
		log.Fatalf("auto-migrate: %v", err)
	}
}
