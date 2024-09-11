package domain

import (
	"context"
	"time"
)

type VcrRequestBatch struct {
	Pkg  uint `json:"package"`
	Size uint `json:"size"`
}

type Vcr struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Pkg       uint      `json:"package"`
	Username  string    `json:"username" gorm:"index"`
	Batchcode string    `json:"batch" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VcrRepository interface {
	First(ctx context.Context, vcr *Vcr) (res Vcr, err error)
	Batch(ctx context.Context, vcrs []Vcr, radcheck []Rdcheck) (err error)
	GetByUsername(ctx context.Context, username string) (vcr Vcr, err error)
	GetByBatchname(batchname string) (vcr []Vcr, err error)
}

type VcrUsecase interface {
	CreateBatch(ctx context.Context, pkg uint, size uint) (batch string, err error)
}
