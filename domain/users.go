package domain

import (
	"context"
	"time"
)

// Users ...
type Users struct {
	ID        *int64  `json:"id"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	Level     *string `json:"level"`
	CreatedAt *string `json:"created_at"`
}

type UsrRequest struct {
	IDLevel  uint   `json:"id_level" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Usr struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	IDLevel   uint      `json:"id_level"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	UserLevel Level     `json:"usr_level,omitempty" gorm:"foreignKey:IDLevel;references:ID"`
	CreatedAt time.Time `json:"created_at"`
}

type UsrResetRequest struct {
	ID uint `json:"id" validate:"required"`
}

type Level struct {
	ID   uint   `json:"id" gorm:"primaryKey" validate:"required"`
	Name string `json:"name" gorm:"unique" validate:"required"`
}

type LevelRequest struct {
	Name string `json:"name" validate:"required"`
}

func (Usr) TableName() string {
	return "users"
}

// UsersRepository ...
type UsersRepository interface {
	First(user *Usr) (res Usr, err error)
	LevelFirst(level *Level) (res Level, err error)
	Save(user *Usr) (err error)
	LevelSave(level *Level) (err error)
	Find(user *Usr) (res []Usr, err error)
	LevelFind(level *Level) (res []Level, err error)
	Delete(user *Usr) (err error)
	LevelDelete(level *Level) (err error)
	// Search(ctx context.Context, user Users) (res []Users, err error)
	// Insert(ctx context.Context, user *Users) (err error)
	// Update(ctx context.Context, user Users) (err error)
}

// UsersUsecase ...
type UsersUsecase interface {
	Save(user *Usr) (err error)
	Find(user *Usr) (res []Usr, err error)
	ResetPassword(request *UsrResetRequest) (err error)
	Delete(request *Usr) (err error)
	Login(c context.Context, username string, password string) (res JWTCustomClaims, err error)
	LevelFirst(level *Level) (res Level, err error)
	LevelFind(level *Level) (res []Level, err error)
	LevelSave(request *LevelRequest) (err error)
	LevelDelete(level *Level) (err error)
}
