package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/middleware"
	"github.com/spf13/viper"
)

type customerHandler struct {
	ucase domain.CustomerUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.CustomerUsecase) {
	handler := &customerHandler{ucase: uc}

	group := e.Group("/api/v1/customers")
	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString(`server.secret`)),
	}))
	group.Use(middleware.AppMiddleware)

	group.GET("", handler.Find)
	group.POST("", handler.Save)
	group.DELETE("", handler.Delete)
}

func (h *customerHandler) Find(e echo.Context) error {
	res, err := h.ucase.Find(context.Background(), 10)
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

func (h *customerHandler) Save(e echo.Context) error {

	req := domain.CustomerRequest{
		Name:      e.FormValue("name"),
		Sex:       e.FormValue("sex"),
		Birthdate: e.FormValue("birth_date"),
		Phone:     e.FormValue("phone"),
		Identity:  e.FormValue("identity"),
		Email:     e.FormValue("email"),
		Job:       e.FormValue("job"),
		Lat:       e.FormValue("lat"),
		Lng:       e.FormValue("lng"),
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		return e.JSON(http.StatusFailedDependency, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := e.Request().ParseMultipartForm(1024); err != nil {
		return e.JSON(http.StatusFailedDependency, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	uploadedFile, handler, err := e.Request().FormFile("image")
	if err != nil {
		return e.JSON(http.StatusFailedDependency, domain.ErrorMessage{
			Success: false,
			Message: "avatar tidak ditemukan",
		})
	}

	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	fmt.Println("handler.Filename", handler.Filename)

	publicPath := "/public/upload/img"

	filename := fmt.Sprintf("%s%s", req.Identity, filepath.Ext(handler.Filename))

	filelocation := filepath.Join(dir, ("." + publicPath), filename)

	targetFile, err := os.OpenFile(filelocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: err.Error(),
		})
	}

	req.Image = fmt.Sprintf("%s/%s", publicPath, filename)

	err = h.ucase.Save(context.Background(), &req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, domain.ErrorMessage{
			Success: false,
			Message: "internal server error",
		})
	}

	return e.JSON(http.StatusOK, domain.ErrorMessage{
		Success: true,
		Message: "berhasil menyimpan data",
	})
}

func (h *customerHandler) Delete(c echo.Context) error {
	request := new(domain.CustomerRequestUpDel)
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

	err = h.ucase.Delete(request)
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
