package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/p3rch1/review-manager/internal/models"
)

func (api *ServiceAPI) Stats(ctx echo.Context) error {
	var (
		resp                       models.StatsResponse
		serviceErr, userErr, prErr error
		wg                         sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp.Service, serviceErr = api.DB.ServiceStats()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp.User, userErr = api.DB.UserStats()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp.PR, prErr = api.DB.PRStats()
	}()

	wg.Wait()
	if err := errors.Join(serviceErr, userErr, prErr); err != nil {
		return fmt.Errorf("stats fetch: %w", err)
	}

	return ctx.JSON(http.StatusOK, &resp)
}
