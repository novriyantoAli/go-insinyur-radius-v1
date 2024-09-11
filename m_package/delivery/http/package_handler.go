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
	"gorm.io/gorm"
)

type packageHandler struct {
	ucase domain.PackageUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.PackageUsecase) {
	handler := &packageHandler{ucase: uc}

	group := e.Group("/api/v1/packages")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("", handler.Find)
	group.GET("/detail", handler.First)
	group.POST("", handler.Save)
	group.PUT("", handler.Update)
	group.DELETE("", handler.Delete)
}

func (h *packageHandler) Find(e echo.Context) error {
	res, err := h.ucase.Find()
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

func (h *packageHandler) First(e echo.Context) error {
	id := e.QueryParam("pkgs")
	u64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusBadRequest, domain.ErrorMessage{
			Success: false,
			Message: "kesalahan parameter request",
		})
	}
	pkg, err := h.ucase.First(uint(u64))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server sedang mengalami kendala, coba lagi nanti",
		})
	}

	return e.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "data berhasil di ambil",
		Data:    pkg,
	})
}

func (h *packageHandler) Save(e echo.Context) error {
	request := new(domain.PkgRequest)
	err := e.Bind(request)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := e.Validate(request); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	// convert request validity value
	validityValue, err := strconv.ParseUint(request.ValidityValue, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	price, err := strconv.ParseUint(request.Price, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	margin, err := strconv.ParseUint(request.Margin, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	err = h.ucase.Create(context.Background(), &domain.Pkg{
		Name:          request.Name,
		ValidityUnit:  request.ValidateUnit,
		Profile:       request.Profile,
		ValidityValue: uint(validityValue),
		Price:         uint(price),
		Margin:        uint(margin),
	})

	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "terjadi kesalahan ketika menyimpan data",
		})
	}

	return e.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: false,
		Message: "berhasil menyimpan data",
	})
}

func (h *packageHandler) Update(c echo.Context) error {
	request := new(domain.PkgRequest)
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

	// convert request to uint
	validityValue, err := strconv.ParseUint(request.ValidityValue, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	price, err := strconv.ParseUint(request.Price, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	margin, err := strconv.ParseUint(request.Margin, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	pkg, err := h.ucase.First(request.ID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return c.JSON(http.StatusNotFound, domain.ErrorMessage{
				Success: false,
				Message: "data tidak ditemukan...",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
				Success: false,
				Message: "hubungi penyedia layanan aplikasi...",
			})
		}
	}

	pkg.Name = request.Name
	pkg.ValidityUnit = request.ValidateUnit
	pkg.Profile = request.Profile
	pkg.Margin = uint(margin)
	pkg.Price = uint(price)
	pkg.ValidityValue = uint(validityValue)

	err = h.ucase.Save(&pkg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "hubungi penyedia layanan aplikasi...",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil menyimpan data",
	})
}

func (h *packageHandler) Delete(c echo.Context) error {
	request := new(domain.PkgRequest)
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

	err = h.ucase.Delete(request.ID)
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
