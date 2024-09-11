package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	Conn *gorm.DB
}

func NewMysqlRepository(conn *gorm.DB) domain.RadgroupRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) FindCheck(rgck *domain.Radgroupcheck) (res []domain.Radgroupcheck, err error) {
	err = m.Conn.Find(&res, rgck).Error
	return
}

func (m *mysqlRepository) FindReply(rgry *domain.Radgroupreply) (res []domain.Radgroupreply, err error) {
	err = m.Conn.Find(&res, rgry).Error
	return
}

func (m *mysqlRepository) SaveCheck(rgck *domain.Radgroupcheck) (err error) {
	err = m.Conn.Create(rgck).Error
	return
}

func (m *mysqlRepository) SaveReply(rgry *domain.Radgroupreply) (err error) {
	err = m.Conn.Create(rgry).Error
	return
}

func (m *mysqlRepository) Save(ctx context.Context, cst *domain.Customer) (err error) {
	err = m.Conn.Save(cst).Error
	return
}

func (m *mysqlRepository) DeleteCheck(rgck *domain.Radgroupcheck) (err error) {
	err = m.Conn.Delete(&rgck).Error
	return
}

func (m *mysqlRepository) DeleteReply(rgry *domain.Radgroupreply) (err error) {
	err = m.Conn.Delete(&rgry).Error
	return
}
