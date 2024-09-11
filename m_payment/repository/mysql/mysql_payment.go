package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	Conn *gorm.DB
}

func NewMysqlRepository(conn *gorm.DB) domain.PaymentRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) Find(ctx context.Context, payment *domain.Payment) (res []domain.Payment, err error) {
	err = m.Conn.Find(&res, payment).Error
	return
}

func (m *mysqlRepository) PaymentOrder(ctx context.Context, payment *domain.Payment) (err error) {
	err = m.Conn.Create(payment).Error
	return
}

func (m *mysqlRepository) Today(ctx context.Context, todaystart string, todayend string) (res []domain.Payment, err error) {
	err = m.Conn.Where("created_at BETWEEN ? AND ?", todaystart, todayend).Find(&res).Error
	return
}

func (m *mysqlRepository) Month(ctx context.Context, targetmonth int, targetyear int) (res []domain.Payment, err error) {
	err = m.Conn.Model(domain.Payment{}).Where("MONTH(created_at) = ? AND YEAR(created_at) = ? ", targetmonth, targetyear).Find(&res).Error
	return
}

func (m *mysqlRepository) Year(ctx context.Context, targetyear int) (res []domain.Payment, err error) {
	err = m.Conn.Model(domain.Payment{}).Where("YEAR(created_at) = ? ", targetyear).Find(&res).Error
	return
}
