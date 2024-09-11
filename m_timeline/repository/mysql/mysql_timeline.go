package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	DB *gorm.DB
}

func NewMysqlRepository(db *gorm.DB) domain.TimelineRepository {
	return &mysqlRepository{DB: db}
}

func (m *mysqlRepository) Find(ctx context.Context) (res []domain.Timeline, err error) {
	err = m.DB.Find(&res, &domain.Timeline{}).Limit(10).Error
	return
}
