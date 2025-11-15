package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func (api *ServiceAPI) SetIsActive(ctx echo.Context) error {
	var req models.SetActiveRequest
	if err := ctx.Bind(&req); err != nil {
		return models.ErrInvalidInput
	}

	user, err := api.DB.SetIsActive(ctx.Request().Context(), &req)
	if err != nil {
		return fmt.Errorf("set is_active DB: %w", err)
	}

	return ctx.JSON(http.StatusOK, map[string]any{"user": user})
}

func (api *ServiceAPI) GetReviews(ctx echo.Context) error {
	id := ctx.QueryParam("user_id")
	if id == "" {
		return models.ErrInvalidInput
	}

	reviews, err := api.DB.GetReviews(ctx.Request().Context(), id)
	if err != nil {
		return fmt.Errorf("get reviews DB: %w", err)
	}

	return ctx.JSON(http.StatusOK, &models.GetReviewResponse{
		UserID:       id,
		PullRequests: reviews,
	})
}
