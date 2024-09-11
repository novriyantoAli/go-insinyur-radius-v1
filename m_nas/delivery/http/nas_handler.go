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

type nasHandler struct {
	ucase domain.NasUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.NasUsecase) {
	handler := &nasHandler{ucase: uc}

	group := e.Group("/api/v1/nas")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("", handler.Find)
	group.POST("", handler.Create)
	group.DELETE("", handler.Delete)
}

func (h *nasHandler) Find(e echo.Context) error {
	res, err := h.ucase.Find(context.Background(), &domain.Nas{})
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

func (h *nasHandler) Create(e echo.Context) error {
	requestModel := new(domain.NasReq)
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

	err = h.ucase.Create(context.Background(), requestModel)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "an internal server error",
		})
	}

	return e.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "nas ditambahkan",
	})
}

func (h *nasHandler) Delete(c echo.Context) error {
	request := new(domain.NasDeleteReq)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "err when bind data",
		})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "err when validate request",
		})
	}

	err = h.ucase.Delete(context.Background(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server error hubungi penyedia layanan",
		})
	}

	return c.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "nas telah dihapus...",
	})
}
