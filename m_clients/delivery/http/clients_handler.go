package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type clientHandler struct {
	ucase domain.ClientUsecase
}

func NewHandler(e *echo.Echo, uc domain.ClientUsecase) {
	handler := &clientHandler{ucase: uc}

	group := e.Group("/api/v1/clients")

	// private route
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("", handler.Fetch)
	group.GET("/binding", handler.FetchBinding)
	// group.GET("/detail", handler.Detail)
	group.POST("", handler.Save)
	group.POST("/binding", handler.SaveBinding)
	group.DELETE("", handler.Delete)
	group.DELETE("/binding", handler.DeleteBinding)
	// group.POST("/subscribe", handler.Subscribe)
}

func (h *clientHandler) Fetch(c echo.Context) error {
	username := c.QueryParam("username")
	request := new(domain.Client)
	request.Username = username

	res, err := h.ucase.Find(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "data berhasil di ambil",
		Data:    res,
	})
}

func (h *clientHandler) FetchBinding(c echo.Context) error {
	mac := c.QueryParam("mac")
	request := new(domain.ClientBinding)
	request.Mac = mac

	res, err := h.ucase.FindBinding(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika mengambil data",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "data berhasil di ambil",
		Data:    res,
	})
}

func (h *clientHandler) Detail(e echo.Context) error {
	id := e.QueryParam("cli")
	u64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusBadRequest, domain.ErrorMessage{
			Success: false,
			Message: "kesalahan parameter request",
		})
	}

	res, err := h.ucase.First(uint(u64))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan pada server",
		})
	}

	return e.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "berhasil mengambil data",
		Data:    res,
	})
}

func (h *clientHandler) Save(c echo.Context) error {
	m := new(domain.ClientRequest)
	err := c.Bind(&m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "permintaan tidak lengkap",
		})
	}

	if err := c.Validate(m); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	priceUInt64, err := strconv.ParseUint(m.Price, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	sessionUInt64, err := strconv.ParseUint(m.Session, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	validityValueUInt64, err := strconv.ParseUint(m.ValidityValue, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}
	client := domain.Client{
		IDCustomer:    m.IDCustomer,
		Username:      m.Username,
		Password:      m.Password,
		Price:         uint(priceUInt64),
		Session:       uint(sessionUInt64),
		Speed:         m.Speed,
		Notes:         m.Notes,
		ValidityUnit:  m.ValidityUnit,
		ValidityValue: uint(validityValueUInt64),
	}

	err = h.ucase.Save(&client)
	if err != nil {
		if errors.Is(gorm.ErrDuplicatedKey, err) {
			return c.JSON(http.StatusConflict, domain.ErrorMessage{
				Success: false,
				Message: "pengguna telah terdaftar sebagai client sebelumnya",
			})
		}
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "kesalahan server, silahkan coba beberapa saat lagi",
		})
	}

	return c.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "data berhasil disimpan",
	})
}

func (h *clientHandler) SaveBinding(c echo.Context) error {
	m := new(domain.ClientBinding)
	err := c.Bind(&m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "permintaan tidak lengkap",
		})
	}

	if err := c.Validate(m); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	err = h.ucase.SaveBinding(m)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "kesalahan server, silahkan coba beberapa saat lagi",
		})
	}

	return c.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "data berhasil disimpan",
	})
}

func (h *clientHandler) Delete(c echo.Context) error {
	request := new(domain.ClientRequest)
	// pkgreq := new(domain.Pkgdeleterequest)
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

	err = h.ucase.Delete(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika menghapus data",
		})
	}

	return c.JSON(http.StatusNoContent, domain.ErrorMessage{
		Success: true,
		Message: "data telah dihapus",
	})
}

func (h *clientHandler) DeleteBinding(c echo.Context) error {
	request := new(domain.ClientBinding)
	// pkgreq := new(domain.Pkgdeleterequest)
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

	err = h.ucase.DeleteBinding(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika menghapus data",
		})
	}

	return c.JSON(http.StatusNoContent, domain.ErrorMessage{
		Success: true,
		Message: "data telah dihapus",
	})
}
