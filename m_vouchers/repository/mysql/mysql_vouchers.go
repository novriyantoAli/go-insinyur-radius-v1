package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	Conn *gorm.DB
}

func NewMysqlRepository(conn *gorm.DB) domain.VcrRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) First(ctx context.Context, voucher *domain.Vcr) (res domain.Vcr, err error) {
	err = m.Conn.First(&res, voucher).Error
	return
}

func (m *mysqlRepository) GetByUsername(ctx context.Context, username string) (vcr domain.Vcr, err error) {
	err = m.Conn.Where(&domain.Vcr{Username: username}).First(&vcr).Error
	return
}

func (m *mysqlRepository) GetByBatchname(batchname string) (vcrs []domain.Vcr, err error) {
	err = m.Conn.Where(&domain.Vcr{Batchcode: batchname}).Find(&vcrs).Error
	return
}

func (m *mysqlRepository) Batch(ctx context.Context, vcrs []domain.Vcr, radcheck []domain.Rdcheck) (err error) {
	tx := m.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(vcrs).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(radcheck).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
