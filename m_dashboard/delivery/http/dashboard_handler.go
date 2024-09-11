package http

import (
	"context"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
)

type dashboardHandler struct {
	ucase domain.DashboardUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.DashboardUsecase) {
	handler := &dashboardHandler{ucase: uc}

	group := e.Group("/api/v1/dashboards")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("", handler.All)
	group.GET("/timeline", handler.AllTimeline)
	group.GET("/tlocation", handler.AllTopLocation)
}

func (hn *dashboardHandler) AllTopLocation(c echo.Context) error {
	res, err := hn.ucase.TopLocation(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    res,
	})
}

func (hn *dashboardHandler) AllTimeline(c echo.Context) error {
	res, err := hn.ucase.Timeline(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    res,
	})
}

func (hn *dashboardHandler) All(c echo.Context) error {
	res, err := hn.ucase.QuickCount(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    res,
	})
}
