package domain

import "context"

type Nas struct {
	ID          uint   `json:"id"`
	Nasname     string `json:"nasname"`
	Shortname   string `json:"shortname"`
	Type        string `json:"type"`
	Ports       uint   `json:"port"`
	Secret      string `json:"secret"`
	Server      string `json:"server"`
	Community   string `json:"community"`
	Description string `json:"description"`
}

type NasReq struct {
	ID          uint   `json:"id" validate:"required"`
	Nasname     string `json:"nasname" validate:"required"`
	Shortname   string `json:"shortname" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Ports       string `json:"port"`
	Secret      string `json:"secret" validate:"required"`
	Server      string `json:"server"`
	Community   string `json:"community"`
	Description string `json:"description"`
}

type NasDeleteReq struct {
	ID        uint   `json:"id" validate:"required"`
	Nasname   string `json:"nasname" validate:"required"`
	Shortname string `json:"shortname" validate:"required"`
	Type      string `json:"type" validate:"required"`
	Secret    string `json:"secret" validate:"required"`
}

func (Nas) TableName() string {
	return "nas"
}

type NasRepository interface {
	Find(ctx context.Context, param *Nas) (res []Nas, err error)
	Create(ctx context.Context, param *Nas) (err error)
	Delete(ctx context.Context, param *Nas) (err error)
}

type NasUsecase interface {
	Find(c context.Context, nas *Nas) (res []Nas, err error)
	Create(ctx context.Context, param *NasReq) (err error)
	Delete(ctx context.Context, param *NasDeleteReq) (err error)
}
