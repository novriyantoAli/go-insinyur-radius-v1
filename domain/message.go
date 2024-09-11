package domain

import "context"

type User struct {
	Username string `json:"username"`
}

type MessageUsecase interface {
	ReadUserMessage(c context.Context, usr *User) (err error)
}
