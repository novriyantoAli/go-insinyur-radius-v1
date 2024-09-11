package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
)

type paymentHandler struct {
	ucase domain.PaymentUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.PaymentUsecase) {
	handler := &paymentHandler{ucase: uc}

	group := e.Group("/api/v1/payments")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("/today/count", handler.TodayCount)
	group.GET("/month/current", handler.CurrentMonthCount)
	group.GET("/year/current", handler.CurrentYearCount)
	group.POST("", handler.Pay)
	// group.POST("/client", handler.Client)
}

func (hn *paymentHandler) Pay(c echo.Context) error {
	request := new(domain.PayReq)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	// convert request validity value
	amountValue, err := strconv.ParseUint(request.Amount, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	err = hn.ucase.PaymentOrder(context.Background(), request.IDBatch, amountValue)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	return c.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "berhasil membayar tagihan",
	})

}

func (hn *paymentHandler) CurrentYearCount(c echo.Context) error {
	res, err := hn.ucase.CurrentYear(context.Background())
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

func (hn *paymentHandler) CurrentMonthCount(c echo.Context) error {
	res, err := hn.ucase.CurrentMonth(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	var credit float64 = 0
	var debit float64 = 0
	for _, payment := range res {
		if payment.Status == "kredit" {
			credit += payment.Amount
		} else {
			debit += payment.Amount
		}
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data: domain.DebCreResp{
			Debt:   debit,
			Credit: credit,
		},
	})
}

func (hn *paymentHandler) TodayCount(c echo.Context) error {
	res, err := hn.ucase.TodayReport(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	var total uint = 0
	for _, payment := range res {
		if payment.Status == "kredit" {
			total += uint(payment.Amount)
		}
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    total,
	})
}
