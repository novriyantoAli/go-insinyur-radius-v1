package mysql

import (
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlClients struct {
	DB *gorm.DB
}

func NewMysqlRepository(gorm *gorm.DB) domain.ClientRepository {
	return &mysqlClients{DB: gorm}
}

func (m *mysqlClients) Find(client *domain.Client) (res []domain.Client, err error) {
	err = m.DB.Model(domain.Client{}).Preload("Customer").Find(&res, client).Error
	return
}

func (m *mysqlClients) FindBinding(binding *domain.ClientBinding) (res []domain.ClientBinding, err error) {
	err = m.DB.Find(&res, binding).Error
	return
}

func (m *mysqlClients) First(client *domain.Client) (res domain.Client, err error) {
	err = m.DB.Model(domain.Client{}).Preload("Customer").Preload("ClientBinding").First(&res, client).Error
	return
}

func (m *mysqlClients) Save(client *domain.Client) (err error) {
	err = m.DB.Save(client).Error
	return
}

func (m *mysqlClients) SaveBinding(binding *domain.ClientBinding) (err error) {
	err = m.DB.Save(binding).Error
	return
}

func (m *mysqlClients) Delete(client *domain.Client) (err error) {
	err = m.DB.Delete(&client).Error
	return
}

func (m *mysqlClients) DeleteBinding(binding *domain.ClientBinding) (err error) {
	err = m.DB.Delete(&binding).Error
	return
}
