package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"taskTrack/internal/config"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(ctx context.Context, config *config.Config, logger *zap.Logger, router *gin.Engine) {
	notifContext, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
		Addr:         config.App.Host + ":" + config.App.Port,
		Handler:      router,
	}

	logger.Info("listen: " + config.App.Host + ":" + config.App.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: " + err.Error())
		}
	}()

	<-notifContext.Done()
	stop()
	logger.Info("Shutting down graceful")

	notifContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(notifContext); err != nil {
		logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Info("Server exiting")
}
