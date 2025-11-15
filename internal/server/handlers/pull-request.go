package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func (api *ServiceAPI) CreatePR(ctx echo.Context) error {
	var req models.PRCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}

	pr, err := api.DB.CreatePR(ctx.Request().Context(), &req, api.Config.Business.ReviewersCount)
	if err != nil {
		return fmt.Errorf("create pr DB: %w", err)
	}

	return ctx.JSON(http.StatusCreated, map[string]any{"pr": pr})
}

func (api *ServiceAPI) MergePR(ctx echo.Context) error {
	var req models.MergeRequest
	if err := ctx.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}

	pr, err := api.DB.Merge(ctx.Request().Context(), &req)
	if err != nil {
		return fmt.Errorf("merge pr DB: %w", err)
	}

	return ctx.JSON(http.StatusOK, map[string]any{"pr": pr})
}

func (api *ServiceAPI) ReassignPR(ctx echo.Context) error {
	var req models.ReassignRequest
	if err := ctx.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}

	resp, err := api.DB.ReassignPR(ctx.Request().Context(), &req)
	if err != nil {
		return fmt.Errorf("reassign pr DB: %w", err)
	}

	return ctx.JSON(http.StatusOK, resp)
}
