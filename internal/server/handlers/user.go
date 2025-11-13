package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func (api *ServiceAPI) SetUserIsActive(ctx echo.Context) error {
	var request models.SetActiveRequest
	if err := ctx.Bind(&request); err != nil {
		return models.ErrInvalidInput
	}

	user, err := api.DB.SetIsActive(ctx.Request().Context(), &request)
	if err != nil {
		return fmt.Errorf("set is_active DB: %w", err)
	}

	return ctx.JSON(http.StatusOK, map[string]any{"user": user})
}

func (api *ServiceAPI) GetUserReviewPRs(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return models.ErrInvalidInput
	}

	return nil
}
