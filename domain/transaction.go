package domain

import (
	"context"
	"time"
)

type Transaction struct {
	ID              *int64     `json:"id"`
	IDReseller      *int64     `json:"id_reseller"`
	IDRadpackage    *int64     `json:"id_radpackage"`
	TransactionCode *string    `json:"transaction_code"`
	Status          *string    `json:"status"`
	Value           *int64     `json:"value"`
	Information     *string    `json:"information"`
	CreatedAt       *time.Time `json:"created_at"`
	Radpackage      Radpackage `json:"radpackage"`
	Reseller        Users      `json:"reseller"`
}

type TransactionPage struct {
	TotalPage int64         `json:"total_page"`
	Data      []Transaction `json:"data"`
}

type TransactionRepository interface {
	CountPage(ctx context.Context, spec Transaction)(res int64, err error)
	Fetch(ctx context.Context, idUsers int64,  id int64, limit int64) (res []Transaction, err error)
	Get(ctx context.Context, transaction Transaction) (res []Transaction, err error)
}

type TransactionUsecase interface {
	Fetch(c context.Context, idUsers int64, id int64, limit int64, query string) (res TransactionPage, err error)
}
