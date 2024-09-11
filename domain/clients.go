package domain

import "time"

type ClientRequest struct {
	ID            uint         `json:"id"`
	IDCustomer    uint         `json:"id_customer" validate:"required"`
	Username      string       `json:"username" validate:"required"`
	Password      string       `json:"password" validate:"required"`
	Price         string       `json:"price" validate:"required"`
	Speed         string       `json:"speed" validate:"required"`
	Session       string       `json:"session" validate:"required"`
	ValidityUnit  ValidityUnit `json:"validity_unit" validate:"required"`
	ValidityValue string       `json:"validity_value" validate:"required"`
	Notes         string       `json:"notes"`
}

type Client struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	IDCustomer    uint            `json:"id_customer" gorm:"index;unique"`
	Username      string          `json:"username" gorm:"index"`
	Password      string          `json:"password"`
	ValidityUnit  ValidityUnit    `gorm:"type:enum('HOUR','DAY','MONTH')" json:"validity_unit"`
	ValidityValue uint            `json:"validity_value"`
	Price         uint            `json:"price"`
	Speed         string          `json:"speed"`
	Session       uint            `json:"session"`
	Notes         string          `json:"notes"`
	ClientBinding []ClientBinding `json:"binding" gorm:"foreignKey:IDClient;references:ID"`
	Customer      Customer        `json:"customer,omitempty" gorm:"foreignKey:IDCustomer;references:ID"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type ClientBinding struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	IDClient uint   `json:"id_client" gorm:"index" validate:"required"`
	Mac      string `json:"mac" validate:"required"`
}

type ClientRepository interface {
	Find(client *Client) (res []Client, err error)
	FindBinding(binding *ClientBinding) (res []ClientBinding, err error)
	First(client *Client) (res Client, err error)
	Save(client *Client) (err error)
	SaveBinding(binding *ClientBinding) (err error)
	Delete(client *Client) (err error)
	DeleteBinding(binding *ClientBinding) (err error)
}

type ClientUsecase interface {
	Find(client *Client) (res []Client, err error)
	FindBinding(binding *ClientBinding) (res []ClientBinding, err error)
	First(id uint) (res Client, err error)
	Save(client *Client) (err error)
	SaveBinding(binding *ClientBinding) (err error)
	Delete(request *ClientRequest) (err error)
	DeleteBinding(binding *ClientBinding) (err error)
}
