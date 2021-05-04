package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested items is not exists
	ErrNotFound = errors.New("Your requested Items is Not Found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your item already exists")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given param is not valid")

	ErrInSufficientBalance = errors.New("in sufficient balance")
	ErrUsernameOrPasswordWrong = errors.New("username or password is not valid")
)
