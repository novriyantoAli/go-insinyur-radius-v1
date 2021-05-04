package domain

import (
	"context"
	"time"
)

type Package struct {
	ID            *int64     `json:"id"`
	Name          *string    `json:"name"`
	ValidityValue *int64     `json:"validity_value"`
	ValidityUnit  *string    `json:"validity_unit"`
	Price         *int64     `json:"price"`
	Margin        *int64     `json:"margin"`
	Profile       *string    `json:"profile"`
	CreatedAt     *time.Time `json:"created_at"`
}

type PackagePage struct {
	TotalPage int64     `json:"total_page"`
	Data      []Package `json:"data"`
}

type PackageRepository interface {
	CountPage(ctx context.Context, spec Package)(res int64, err error)
	Fetch(ctx context.Context, id int64, limit int64) (res []Package, err error)
	Get(ctx context.Context, packages Package) (res []Package, err error)
}

type PackageUsecase interface {
	Fetch(c context.Context, id int64, limit int64) (res PackagePage, err error)
}
