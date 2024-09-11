package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
)

type orderHandler struct {
	ucase domain.OrderUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.OrderUsecase) {
	handler := &orderHandler{ucase: uc}

	group := e.Group("/api/v1/orders")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("", handler.Find)
	group.POST("", handler.Order)
	group.POST("/client", handler.Client)
	group.GET("/report/hotspot/month", handler.ReportHotspotMonth)
	group.GET("/report/client/month", handler.ReportClientMonth)
}

func (h *orderHandler) Find(e echo.Context) error {
	last := e.QueryParam("last")
	lastID, err := strconv.ParseUint(last, 10, 64)
	if err != nil {
		lastID = 0
	}

	limit := e.QueryParam("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	res, err := h.ucase.Find(context.Background(), uint(lastID), limitInt)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	return e.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    res,
	})
}

func (hn *orderHandler) ReportClientMonth(c echo.Context) error {
	client, err := hn.ucase.ReportClientMonth(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan pada server...",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    client,
	})
}

func (hn *orderHandler) ReportHotspotMonth(c echo.Context) error {
	hotspot, err := hn.ucase.ReportHotspotMonth(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan pada server...",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    hotspot,
	})
}

func (hn *orderHandler) Client(e echo.Context) error {
	requestModel := new(domain.OrderClientRequest)
	err := e.Bind(requestModel)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := e.Validate(requestModel); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	err = hn.ucase.Client(context.Background(), requestModel)
	if err != nil {
		if errors.Is(err, domain.ErrMinimumAmountRequired) {
			return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
				Success: false,
				Message: err.Error(),
			})
		}
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "an internal server error",
		})
	}

	return e.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "batch created...",
	})
}

func (hn *orderHandler) Order(e echo.Context) error {
	requestModel := new(domain.OrderRequest)
	err := e.Bind(requestModel)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := e.Validate(requestModel); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	fmt.Println(requestModel)

	err = hn.ucase.Order(context.Background(), requestModel)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "an internal server error",
		})
	}

	return e.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "batch created...",
	})
}
