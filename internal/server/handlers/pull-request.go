package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func (api *ServiceAPI) CreatePullRequest(c echo.Context) error {
	var req models.PullRequestCreateRequest
	if err := c.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}
	return nil
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
