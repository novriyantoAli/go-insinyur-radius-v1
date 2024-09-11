package http

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/helper"
	"github.com/novriyantoAli/go-insinyur-radius-v1/initializers"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ResponseError ...
type ResponseError struct {
	Message string `json:"error"`
}

type login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type usersHandler struct {
	ucase domain.UsersUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.UsersUsecase) {
	handler := &usersHandler{ucase: uc}

	// isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: []byte(viper.GetString(`server.secret`)),
	// })

	group := e.Group("/api/v1/")

	// route
	group.POST("login", handler.Login)
	// private route
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.POST("reset", handler.ResetPassword)
	group.POST("register", handler.Register)
	group.GET("users", handler.Fetch)
	group.DELETE("users", handler.Delete)

	group.POST("user/reset", handler.ResetPassword)

	levels := group.Group("user/levels")
	levels.GET("", handler.FetchLevel)
	levels.POST("", handler.InsertLevel)
	levels.DELETE("", handler.DeleteLevel)
}

func (hn *usersHandler) FetchLevel(e echo.Context) error {
	res, err := hn.ucase.LevelFind(&domain.Level{})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "",
		})
	}
	return e.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "fetched...",
		Data:    res,
	})
}

func (hn *usersHandler) InsertLevel(e echo.Context) error {
	request := new(domain.LevelRequest)
	err := e.Bind(request)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "err when bind data",
		})
	}
	if err := e.Validate(request); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "err when validate request",
		})
	}

	err = hn.ucase.LevelSave(request)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server sedang error hubungi penyedia layanan",
		})
	}

	return e.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "level created",
	})
}

func (hn *usersHandler) DeleteLevel(c echo.Context) error {
	request := new(domain.Level)
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

	err = hn.ucase.LevelDelete(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server error hubungi penyedia layanan",
		})
	}

	return c.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "level telah dihapus...",
	})
}

func (hn *usersHandler) Fetch(e echo.Context) error {
	res, err := hn.ucase.Find(&domain.Usr{})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "error ehwn fetch data",
		})
	}

	return e.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "fetched...",
		Data:    res,
	})
}

func (hn *usersHandler) ResetPassword(c echo.Context) error {
	request := new(domain.UsrResetRequest)
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
			Message: err.Error(),
		})
	}

	err = hn.ucase.ResetPassword(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "server error hubungi penyedia layanan",
		})
	}

	return c.JSON(http.StatusAccepted, domain.ErrorMessage{
		Success: true,
		Message: "pengguna telah dihapus...",
	})
}

func (hn *usersHandler) Delete(c echo.Context) error {
	request := new(domain.Usr)
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
		Message: "pengguna telah dihapus...",
	})
}

func (hn *usersHandler) Register(e echo.Context) error {
	request := new(domain.UsrRequest)
	err := e.Bind(request)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "err when bind data",
		})
	}

	if err := e.Validate(request); err != nil {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "err when validate request",
		})
	}

	// lakukan pemeriksaan dengan id level
	level, err := hn.ucase.LevelFirst(&domain.Level{ID: request.IDLevel})
	if err != nil {
		return e.JSON(http.StatusBadRequest, domain.ErrorMessage{
			Success: false,
			Message: "level pengguna tidak dikenal...",
		})
	}

	if !initializers.ValidateAAA(level.Name) {
		return e.JSON(http.StatusUnprocessableEntity, domain.ErrorMessage{
			Success: false,
			Message: "level tidak terdaftar",
		})
	}

	usr := domain.Usr{
		IDLevel:  level.ID,
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	err = hn.ucase.Save(&usr)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "error when try to save user",
		})
	}

	return e.JSON(http.StatusCreated, domain.ErrorMessage{
		Success: true,
		Message: "user created",
		Data:    usr,
	})
}

// Login ...
func (hn *usersHandler) Login(e echo.Context) error {
	// get query param
	u := new(login)
	err := e.Bind(u)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusFailedDependency, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := e.Validate(u); err != nil {
		return e.JSON(http.StatusFailedDependency, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	res, err := hn.ucase.Login(e.Request().Context(), u.Username, u.Password)
	if err != nil {
		logrus.Error(err)
		return e.JSON(helper.TranslateError(err), domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	return e.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "successfully to login...",
		Data:    res,
	})
}
