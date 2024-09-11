package domain

import (
	"context"
	"time"
)

type ValidityUnit string

const (
	Hour  ValidityUnit = "HOUR"
	Day   ValidityUnit = "DAY"
	Month ValidityUnit = "MONTH"
)

type PkgRequest struct {
	ID            uint         `json:"id"`
	Name          string       `json:"name" validate:"required"`
	ValidityValue string       `json:"validity_value" validate:"required"`
	ValidateUnit  ValidityUnit `json:"validity_unit" validate:"required"`
	Price         string       `json:"price" validate:"required"`
	Margin        string       `json:"margin" validate:"required"`
	Profile       string       `json:"profile" validate:"required"`
}

type Pkg struct {
	ID               uint         `gorm:"primaryKey" json:"id"`
	Name             string       `json:"name"`
	ValidityValue    uint         `json:"validity_value"`
	ValidityUnit     ValidityUnit `gorm:"type:enum('HOUR','DAY','MONTH')" json:"validity_unit"`
	Price            uint         `json:"price"`
	Margin           uint         `json:"margin"`
	Profile          string       `json:"profile" gorm:"index"`
	SelectionProfile Radusergroup `json:"selection_profile,omitempty" gorm:"foreignKey:Profile;references:Username"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

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
	Find(pkg *Pkg) (res []Pkg, err error)
	Save(pk *Pkg) (err error)
	Create(ctx context.Context, pkg *Pkg) (err error)
	// CountPage(ctx context.Context, spec Package) (res int64, err error)
	// Fetch(ctx context.Context, id int64, limit int64) (res []Package, err error)
	// Get(ctx context.Context, packages Package) (res []Package, err error)
	First(id uint) (pkg Pkg, err error)
	Delete(pkgs []Pkg) (err error)
}

type PackageUsecase interface {
	Find() (res []Pkg, err error)
	First(id uint) (res Pkg, err error)
	Save(pkg *Pkg) (err error)
	Create(c context.Context, pkg *Pkg) (err error)
	Delete(id uint) (err error)
	// Fetch(c context.Context, id int64, limit int64) (res PackagePage, err error)
}
