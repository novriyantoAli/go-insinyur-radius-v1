package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	Conn *gorm.DB
}

func NewMysqlRepository(conn *gorm.DB) domain.NasRepository {
	return &mysqlRepository{Conn: conn}
}

func (my *mysqlRepository) Find(ctx context.Context, param *domain.Nas) (res []domain.Nas, err error) {
	err = my.Conn.Find(&res, param).Error
	return
}

func (my *mysqlRepository) Create(ctx context.Context, param *domain.Nas) (err error) {
	err = my.Conn.Create(param).Error
	return
}

func (m *mysqlRepository) Delete(ctx context.Context, param *domain.Nas) (err error) {
	err = m.Conn.Delete(&param).Error
	return
}
