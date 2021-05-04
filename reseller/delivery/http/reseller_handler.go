package http

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ResponseError ...
type ResponseError struct {
	Message string `json:"error"`
}

type transaction struct {
	id int `json:"id_package" validate:"required"`
}

// type login struct {
// 	Username string `json:"username" validate:"required"`
// 	Password string `json:"password" validate:"required"`
// }

type resellerHandler struct {
	ucase domain.ResellerUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.ResellerUsecase) {
	handler := &resellerHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	})

	group := e.Group("/api/reseller", isLoggedIn)
	group.POST("", handler.PostTransaction)
	group.GET("", handler.GetBalance)
	group.POST("/change-package", handler.PostChangePackage)
	// group.POST("", handler.Save)
	// group.POST("/find", handler.Find)
	// group.PUT("", handler.Update)
	// group.DELETE("/:id", handler.Delete)
}

func (h *resellerHandler) GetBalance(e echo.Context) error {

	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	res, err := h.ucase.Balance(e.Request().Context(), int64(claims["id"].(float64)))
	if err != nil {
		logrus.Error(err)
		return e.JSON(helper.TranslateError(err), ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

func (hn *resellerHandler) PostTransaction(e echo.Context) error {
	// get query param
	// u := new(transaction)
	// err := e.Bind(u)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	// }

	// if err := e.Validate(u); err != nil {
	// 	return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	// }

	idString := e.FormValue("idPackage")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		id = 0
	}
	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	res, err := hn.ucase.Transaction(e.Request().Context(), int64(claims["id"].(float64)), id)
	if err != nil {
		logrus.Error(err)
		return e.JSON(helper.TranslateError(err), ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

func (h *resellerHandler) PostChangePackage(e echo.Context) error {
	voucher := e.FormValue("voucher")
	profileName := e.FormValue("profile")

	if voucher == "" {
		logrus.Warning("input parameter not valid")
		return e.JSON(http.StatusFailedDependency, helper.ResponseErrorMessage{Message: "input parameter not valid"})
	}

	err := h.ucase.ChangePackage(e.Request().Context(), voucher, profileName)
	if err != nil {
		logrus.Error(err)
		return e.JSON(helper.TranslateError(err), ResponseError{Message: err.Error()})
	}
	return e.JSON(http.StatusOK, helper.ResponseSuccessMessage{Message: "success to change package"})
}
