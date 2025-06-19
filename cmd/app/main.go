package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"taskTrack/internal/config"
	v1 "taskTrack/internal/controllers/http/v1"
	repos1 "taskTrack/internal/services/task/repository"
	usecase1 "taskTrack/internal/services/task/usecase"
	"taskTrack/internal/taskprocessor"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	logger, _ := zap.NewProduction()
	gin.SetMode(cfg.App.Mode)

	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	repository1 := repos1.NewTaskRepository()
	usecase1 := usecase1.NewTaskTrackUsecase(repository1)
	processor := taskprocessor.NewProcessor(repository1, 5)
	processor.Start(ctx)
	defer processor.Stop()
	api := router.Group("/api")
	{
		v1.NewTaskHandler(cfg, api, usecase1)
	}
	Run(ctx, cfg, logger, router)
}
