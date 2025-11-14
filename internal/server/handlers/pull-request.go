package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func (api *ServiceAPI) CreatePullRequest(ctx echo.Context) error {
	var req models.PRCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}

	pr, err := api.DB.PRCreate(ctx.Request().Context(), &req, api.Config.Business.ReviewersCount)
	if err != nil {
		return fmt.Errorf("create pr: %w", err)
	}

	return ctx.JSON(http.StatusCreated, pr)
}

func (api *ServiceAPI) MergePullRequest(c echo.Context) error {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := c.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}
	return nil
}

func (api *ServiceAPI) ReassignPullRequest(c echo.Context) error {
	var req models.ReassignRequest
	if err := c.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}
	return nil
}
