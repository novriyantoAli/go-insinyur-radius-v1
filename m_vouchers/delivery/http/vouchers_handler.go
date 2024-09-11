package http

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
)

type vouchersHandler struct {
	ucase domain.VcrUsecase
}

func NewHandler(e *echo.Echo, uc domain.VcrUsecase) {
	handler := &vouchersHandler{ucase: uc}

	group := e.Group("/api/v1/vouchers")
	// private route

	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.POST("/batch", handler.Batch)
}

func (hn *vouchersHandler) Batch(e echo.Context) error {
	requestModel := new(domain.VcrRequestBatch)
	err := e.Bind(requestModel)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "failed to process input parameters",
		})
	}

	if err := e.Validate(requestModel); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "failed to validate input parameters",
		})
	}

	bacth, err := hn.ucase.CreateBatch(e.Request().Context(), requestModel.Pkg, requestModel.Size)
	if err != nil || bacth == "" {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "an internal server error",
		})
	}

	return e.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "batch created...",
		Data:    bacth,
	})
}
