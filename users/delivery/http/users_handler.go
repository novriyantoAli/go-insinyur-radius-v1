package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/helper"
	"github.com/sirupsen/logrus"
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

	e.POST("/api/login", handler.Login)

	// group := e.Group("/api/users", isLoggedIn)
	// group.GET("", handler.Fetch)
	// group.POST("", handler.Save)
	// group.POST("/find", handler.Find)
	// group.PUT("", handler.Update)
	// group.DELETE("/:id", handler.Delete)
}

// Login ...
func (hn *usersHandler) Login(e echo.Context) error {
	// get query param
	u := new(login)
	err := e.Bind(u)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	if err := e.Validate(u); err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	res, err := hn.ucase.Login(e.Request().Context(), u.Username, u.Password)
	if err != nil {
		logrus.Error(err)
		return e.JSON(helper.TranslateError(err), ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}
