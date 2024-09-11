package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	Conn *gorm.DB
}

func NewMysqlRepository(conn *gorm.DB) domain.CustomerRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) Find(ctx context.Context, customer *domain.Customer) (customers []domain.Customer, err error) {
	err = m.Conn.Find(&customers, customer).Error
	return
}

func (m *mysqlRepository) First(ctx context.Context, customer *domain.Customer) (res domain.Customer, err error) {
	err = m.Conn.First(&res, customer).Error
	return
}

func (m *mysqlRepository) Save(ctx context.Context, cst *domain.Customer) (err error) {
	err = m.Conn.Save(cst).Error
	return
}

func (m *mysqlRepository) Delete(cst *domain.Customer) (err error) {
	err = m.Conn.Delete(&cst).Error
	return
}
