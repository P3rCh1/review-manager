package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func ErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	var handled *models.ErrorResponce
	if errors.As(err, &handled) {
		ctx.JSON(
			handled.Status,
			&struct {
				Error *models.ErrorResponce `json:"error"`
			}{handled},
		)
		return
	}

	ctx.JSON(http.StatusInternalServerError, models.ErrInternal)
}
