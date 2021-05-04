package helper

import (
	"net/http"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
)

type ResponseErrorMessage struct {
	Message string `json:"error"`
}

type ResponseSuccessMessage struct {
	Message string `json:"message"`
}

// TranslateError ...
func TranslateError(err error) int {
	switch err {
	case domain.ErrBadParamInput:
		return http.StatusBadRequest
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
