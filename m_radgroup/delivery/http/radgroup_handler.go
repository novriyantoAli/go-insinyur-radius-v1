package http

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
)

type radgroupHandler struct {
	ucase domain.RadgroupUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.RadgroupUsecase) {
	handler := &radgroupHandler{ucase: uc}

	group := e.Group("/api/v1/radgroup")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("/check", handler.FindCheck)
	group.GET("/reply", handler.FindReply)
	group.POST("/check", handler.SaveCheck)
	group.POST("/reply", handler.SaveReply)
	group.DELETE("/check", handler.DeleteCheck)
	group.DELETE("/reply", handler.DeleteReply)
	// group.POST("", handler.Save)
	// group.DELETE("", handler.Delete)
}

func (h *radgroupHandler) FindCheck(c echo.Context) error {
	groupname := c.QueryParam("groupname")

	radgroupcheck := domain.Radgroupcheck{}
	if groupname != "" {
		radgroupcheck.Groupname = groupname
	}

	res, err := h.ucase.FindCheck(&radgroupcheck)
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

func (h *radgroupHandler) FindReply(c echo.Context) error {
	groupname := c.QueryParam("groupname")

	radgroupreply := domain.Radgroupreply{}
	if groupname != "" {
		radgroupreply.Groupname = groupname
	}

	res, err := h.ucase.FindReply(&radgroupreply)
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

func (h *radgroupHandler) SaveCheck(c echo.Context) error {
	request := new(domain.Radgroupcheck)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "failed to process input parameters",
		})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "failed to validate input parameters",
		})
	}

	err = h.ucase.SaveCheck(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "internal server error...",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "data properti berhasil disimpan...",
	})
}

func (h *radgroupHandler) SaveReply(c echo.Context) error {
	request := new(domain.Radgroupreply)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "failed to process input parameters",
		})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "failed to validate input parameters",
		})
	}

	err = h.ucase.SaveReply(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "internal server error...",
		})
	}

	return c.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "data properti berhasil disimpan...",
	})
}

func (hn *radgroupHandler) DeleteCheck(c echo.Context) error {
	request := new(domain.Radgroupcheck)
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

	err = hn.ucase.DeleteCheck(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server error hubungi penyedia layanan",
		})
	}

	return c.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "properti telah dihapus...",
	})
}

func (hn *radgroupHandler) DeleteReply(c echo.Context) error {
	request := new(domain.Radgroupreply)
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

	err = hn.ucase.DeleteReply(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server error hubungi penyedia layanan",
		})
	}

	return c.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "properti telah dihapus...",
	})
}
