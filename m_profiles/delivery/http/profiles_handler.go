package http

import (
	"fmt"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
)

type profilesHandler struct {
	ucase domain.ProfilesUsecase
}

func NewHandler(e *echo.Echo, ucase domain.ProfilesUsecase) {
	handler := &profilesHandler{ucase: ucase}

	group := e.Group("/api/v1/profiles")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("", handler.Fetch)
	group.POST("", handler.Save)
	group.DELETE("", handler.Delete)
	// g{
	// SigningKey: []byte(viper.GetString(`server.secret`)),
	// e.GET("/api/profiles", handler.Fetch)
}

func (h *profilesHandler) Fetch(e echo.Context) error {
	radusergroup, err := h.ucase.Fetch(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	return e.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "fetched...",
		Data:    radusergroup,
	})
}

func (h *profilesHandler) Save(c echo.Context) error {
	u := new(domain.Radusergrouprequest)
	err := c.Bind(u)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "permintaan tidak lengkap",
		})
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: fmt.Sprintf("validasi gagal: %s", err.Error()),
		})
	}

	err = h.ucase.Save(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "kesalahan menyimpan data..",
		})
	}

	return c.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "profile berhasil dibuat",
	})
}

func (hn *profilesHandler) Delete(c echo.Context) error {
	request := new(domain.Radusergrouprequest)
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

	err = hn.ucase.Delete(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server error hubungi penyedia layanan",
		})
	}

	return c.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "profil telah dihapus...",
	})
}
