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
	idUser    int64 `json:"id_user" validate:"required"`
	idPackage int64 `json:"id_package" validate:"required"`
}

type price struct {
	name     string `json:"name"`
	validity string `json:"validity"`
	price    string `json:"price"`
}

type packageHandler struct {
	ucase domain.PackageUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.PackageUsecase) {
	handler := &packageHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	})

	e.GET("/price", handler.FetchPublic)

	group := e.Group("/api/package", isLoggedIn)
	group.GET("", handler.Fetch)
}

func (h *packageHandler) FetchPublic(e echo.Context) error {
	res, err := h.ucase.Fetch(e.Request().Context(), 0, 10)
	if err != nil {
		logrus.Error(err)
		return e.JSON(helper.TranslateError(err), ResponseError{Message: err.Error()})
	}

	prArr := make([]price, 0)
	for _, d := range res.Data {
		pr := price{}
		pr.name = *d.Name
		pr.validity = strconv.FormatInt(*d.ValidityValue, 10) + *d.ValidityUnit
		pr.price = strconv.FormatInt((*d.Price + *d.Margin), 10)
		prArr = append(prArr, pr)
	}

	return e.JSON(http.StatusOK, prArr)
}

func (h *packageHandler) Fetch(e echo.Context) error {

	idString := e.QueryParam("id")
	limitString := e.QueryParam("limit")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		id = 0
	}

	limit, err := strconv.ParseInt(limitString, 10, 64)
	if err != nil {
		logrus.Error(err)
		limit = 10
	}

	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["level"] != "admin" && claims["level"] != "user" {
		logrus.Error("level unknown")
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: "unknown level"})
	}

	res, err := h.ucase.Fetch(e.Request().Context(), id, limit)
	if err != nil {
		logrus.Error(err)
		return e.JSON(helper.TranslateError(err), ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}
