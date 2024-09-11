package domain

import (
	"context"
	"time"
)

type Timeline struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Information string    `json:"information"`
	CreatedAt   time.Time `json:"created"`
}

type TimelineRepository interface {
	Find(ctx context.Context) (res []Timeline, err error)
}
