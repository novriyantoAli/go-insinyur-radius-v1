package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
)

func AppMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
				Success: false,
				Message: "missing jwt token...",
			})
		}

		if !token.Valid {
			return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
				Success: false,
				Message: "unvalid jwt token...",
			})
		}
		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok {
		// 	return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
		// 		Success: false,
		// 		Message: "failed to cast token",
		// 	})
		// }

		// err := initializers.ENFORCER.LoadPolicy()
		// if err != nil {
		// 	return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
		// 		Success: false,
		// 		Message: "failed to load app policy...",
		// 	})
		// }

		// ok, err = initializers.ENFORCER.Enforce(fmt.Sprint(claims["id"]), c.Request().URL.Path, c.Request().Method)
		// if err != nil {

		// 	logrus.Error(err)

		// 	return c.JSON(http.StatusInternalServerError, domain.ErrorMessage{
		// 		Success: false,
		// 		Message: "gagal ketika mengautorisasi pengguna..",
		// 	})
		// }

		// if !ok {
		// 	return c.JSON(http.StatusForbidden, domain.ErrorMessage{
		// 		Success: false,
		// 		Message: "anda tidak diizinkan untuk mengakses halaman ini...",
		// 	})
		// }

		return next(c)
	}
}
