package domain

import (
	"context"
	"time"
)

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Identity  string    `json:"identity"`
	Sex       string    `json:"sex"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Job       string    `json:"job"`
	Image     string    `json:"image"`
	Birthdate time.Time `json:"birth_date"`
	Lat       string    `json:"lat" validate:"required"`
	Lng       string    `json:"lng" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomerRequest struct {
	Name      string `json:"name" validate:"required"`
	Identity  string `json:"identity" validate:"required"`
	Sex       string `json:"sex" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Job       string `json:"job" validate:"required"`
	Image     string `json:"image"`
	Birthdate string `json:"birth_date" validate:"required"`
	Lat       string `json:"lat" validate:"required"`
	Lng       string `json:"lng" validate:"required"`
}

type CustomerRequestUpDel struct {
	ID        uint   `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Identity  string `json:"identity" validate:"required"`
	Sex       string `json:"sex" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Job       string `json:"job" validate:"required"`
	Image     string `json:"image"`
	Birthdate string `json:"birth_date" validate:"required"`
	Lat       string `json:"lat" validate:"required"`
	Lng       string `json:"lng" validate:"required"`
}

type CustomerRepository interface {
	Find(ctx context.Context, customer *Customer) (res []Customer, err error)
	First(ctx context.Context, customer *Customer) (res Customer, err error)
	Save(ctx context.Context, customer *Customer) (err error)
	Delete(customer *Customer) (err error)
}

type CustomerUsecase interface {
	Find(c context.Context, limit uint) (res []Customer, err error)
	Save(c context.Context, cr *CustomerRequest) (err error)
	Delete(crud *CustomerRequestUpDel) (err error)
}
