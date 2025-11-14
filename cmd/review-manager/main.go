package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/config"
	"github.com/p3rch1/review-manager/internal/logger"
	"github.com/p3rch1/review-manager/internal/server/handlers"
	"github.com/p3rch1/review-manager/internal/server/middleware"
	"github.com/p3rch1/review-manager/internal/storage/postgres"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config parse fail: %s\n", err)
		os.Exit(1)
	}

	logger, err := logger.Setup(&cfg.Logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup logger fail: %s\n", err)
		os.Exit(1)
	}

	db, err := postgres.NewReviewAPI(&cfg.Postgres)
	if err != nil {
		logger.Error(
			"postgres connection fail",
			"error", err,
		)
		os.Exit(1)
	}
	defer db.Close()

	api := handlers.NewServiceAPI(logger, cfg, db)

	router := SetupServer(api)

	go func() {
		logger.Info("start server")
		err := router.Start(router.Server.Addr)
		if err != nil && err != http.ErrServerClosed {
			logger.Error(
				"listen",
				"error", err,
			)

			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("starting shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		logger.Error(
			"shutdown server",
			"error", err,
		)

		return
	}

	logger.Info("server stopped gracefully")
}

func SetupServer(api *handlers.ServiceAPI) *echo.Echo {
	router := echo.New()

	router.Debug = false
	router.HideBanner = true
	router.Logger.SetOutput(io.Discard)

	router.Use(middleware.Recover(api.Logger))
	router.Use(middleware.Logger(api.Logger))

	router.POST("/team/add", api.AddTeam)
	router.GET("/team/get", api.GetTeam)
	router.POST("/users/setIsActive", api.SetUserIsActive)
	router.GET("/users/getReview", api.GetUserReviewPRs)
	router.POST("/pullRequest/create", api.CreatePR)
	router.POST("/pullRequest/merge", api.MergePR)
	router.POST("/pullRequest/reassign", api.ReassignPR)

	router.HTTPErrorHandler = handlers.ErrorHandler

	router.Server.Addr = api.Config.HTTP.Host + ":" + api.Config.HTTP.Port
	router.Server.ReadTimeout = api.Config.HTTP.ReadTimeout
	router.Server.WriteTimeout = api.Config.HTTP.WriteTimeout
	router.Server.IdleTimeout = api.Config.HTTP.IdleTimeout

	return router
}
