package handlers

import (
	"log/slog"

	"github.com/p3rch1/review-manager/internal/config"
	"github.com/p3rch1/review-manager/internal/storage/postgres"
)

type ServiceAPI struct {
	Logger *slog.Logger
	Config *config.Config
	DB     postgres.ReviewAPI
}

func NewServiceAPI(logger *slog.Logger, config *config.Config, db postgres.ReviewAPI) *ServiceAPI {
	return &ServiceAPI{
		Logger: logger,
		Config: config,
		DB:     db,
	}
}
