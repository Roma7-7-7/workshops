package main

import (
	"context"
	"github.com/Roma7-7-7/workshops/wallet/internal/config"
	"github.com/Roma7-7-7/workshops/wallet/internal/middleware/auth"
	"github.com/Roma7-7-7/workshops/wallet/internal/repository/postgre"
	appHttp "github.com/Roma7-7-7/workshops/wallet/internal/server/http"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/validator"
	"github.com/Roma7-7-7/workshops/wallet/internal/services/wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimout = 5 * time.Second

func main() {
	if logger, err := zap.NewDevelopment(); err != nil {
		panic(err)
	} else {
		zap.ReplaceGlobals(logger)
	}

	cfg := config.GetConfig()
	repo := postgre.NewRepository(cfg.DB.DSN)
	aut := auth.NewMiddleware(repo, cfg.JWT)
	service := wallet.NewService(repo)
	server := appHttp.NewServer(service, &validator.Service{}, aut)

	app := gin.Default()

	server.Register(app)

	appServer := &http.Server{
		Addr:    ":5000",
		Handler: app,
	}

	go func() {
		if err := appServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen app", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown Servers ...")
	appCtx, appCancel := context.WithTimeout(context.Background(), shutdownTimout)
	defer appCancel()

	appShutdownFailed := true
	shutdownComplete := make(chan struct{})

	go func() {
		if err := appServer.Shutdown(appCtx); err != nil {
			zap.L().Error("shutdown app", zap.Error(err))
		} else {
			appShutdownFailed = false
			close(shutdownComplete)
		}
	}()

	select {
	case <-time.After(shutdownTimout):
	case <-shutdownComplete:
	}

	if appShutdownFailed {
		zap.L().Info("app server shutdown failed")
		os.Exit(1)
	}

	zap.L().Info("Server Exited")
}
