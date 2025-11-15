package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func ErrorHandler(logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		var handled *models.ErrorResponce
		if errors.As(err, &handled) {
			err := ctx.JSON(
				handled.Status,
				&struct {
					Error *models.ErrorResponce `json:"error"`
				}{handled},
			)
			if err != nil {
				logger.Error("responce with error fail", "error", err)
			}

			return
		}

		err = ctx.JSON(http.StatusInternalServerError, models.ErrInternal)
		if err != nil {
			logger.Error("responce with error fail", "error", err)
		}
	}
}
