package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func (api *ServiceAPI) AddTeam(ctx echo.Context) error {
	var team models.Team
	if err := ctx.Bind(&team); err != nil || team.TeamName == "" {
		return models.ErrInvalidInput
	}

	if err := api.DB.AddTeam(ctx.Request().Context(), &team); err != nil {
		return fmt.Errorf("add team to DB: %w", err)
	}

	return ctx.JSON(http.StatusCreated, map[string]any{"team": team})
}

func (api *ServiceAPI) GetTeam(ctx echo.Context) error {
	name := ctx.QueryParam("team_name")
	if name == "" {
		return models.ErrInvalidInput
	}

	team, err := api.DB.GetTeam(ctx.Request().Context(), name)
	if err != nil {
		return fmt.Errorf("get team from DB: %w", err)
	}

	return ctx.JSON(http.StatusOK, team)
}
