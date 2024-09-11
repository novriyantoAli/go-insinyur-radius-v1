package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	Conn *gorm.DB
}

func NewMysqlRepository(conn *gorm.DB) domain.OrderRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) Find(ctx context.Context, order *domain.Order) (orders []domain.Order, err error) {
	err = m.Conn.Find(&orders, order).Error
	return
}

func (m *mysqlRepository) FindOrderClientExpire(ctx context.Context) (res []domain.OrderClient, err error) {
	// _ := `SELECT * FROM order_clients WHERE expiration = 'Expiration' AND STR_TO_DATE(value, "%d %b %Y") <= CURDATE()`
	err = m.Conn.Model(domain.OrderClient{}).Preload("Client").Preload("Client.ClientBinding").Where(`STR_TO_DATE(expire, "%d %b %Y") <= CURDATE()`).Find(&res).Error
	return
}

func (m *mysqlRepository) FindOrderClient(ctx context.Context, oc *domain.OrderClient) (res []domain.OrderClient, err error) {
	err = m.Conn.Model(domain.OrderClient{}).Preload("Client").Preload("Client.ClientBinding").Where(`STR_TO_DATE(expire, "%d %b %Y") >= CURDATE() AND id_client = ?`, oc.IDClient).Order("expire asc").Find(&res).Error
	return
}

func (m *mysqlRepository) FindMonthCreated(ctx context.Context, targetmonth int, targetyear int) (res []domain.Order, err error) {
	err = m.Conn.Model(domain.Order{}).Preload("User").Preload("Customer").Preload("OrderHotspot").Preload("OrderHotspot.Package").Preload("OrderClient").Where("MONTH(created_at) = ? AND YEAR(created_at) = ? ", targetmonth, targetyear).Find(&res).Error
	return
}

func (m *mysqlRepository) Paginate(ctx context.Context, lastID uint, limit int) (orders []domain.Order, err error) {
	if lastID == 0 {
		err = m.Conn.Model(domain.Order{}).Preload("User").Preload("Payment").Preload("Customer").Preload("OrderHotspot").Preload("OrderHotspot.Package").Preload("OrderHotspot.Vouchers").Preload("OrderClient").Preload("OrderClient.Client").Order("id desc").Limit(limit).Find(&orders, "id > ?", lastID).Error
	} else {
		err = m.Conn.Model(domain.Order{}).Preload("User").Preload("Payment").Preload("Customer").Preload("OrderHotspot").Preload("OrderHotspot.Package").Preload("OrderHotspot.Vouchers").Preload("OrderClient").Preload("OrderClient.Client").Order("id desc").Limit(limit).Find(&orders, "id < ?", lastID).Error
	}
	return
}

func (m *mysqlRepository) First(ctx context.Context, order *domain.Order) (res domain.Order, err error) {
	err = m.Conn.Model(domain.Order{}).Preload("User").Preload("Customer").Preload("OrderHotspot").Preload("OrderHotspot.Vouchers").Preload("OrderClient").Order("id desc").First(&res, "id_batch = ?", order.IDBatch).Error
	// err = m.Conn.First(&res, order).Error
	return
}

func (m *mysqlRepository) Client(ctx context.Context, order *domain.Order, orderClient *domain.OrderClient, radcheck []domain.Rdcheck, radreply []domain.Radreply, payments []domain.Payment) (err error) {
	tx := m.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(orderClient).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(radcheck).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(radreply).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(payments).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (m *mysqlRepository) Batch(ctx context.Context, order *domain.Order, hotspot *domain.OrderHotspot, vcrs []domain.Vcr, radcheck []domain.Rdcheck, payments []domain.Payment) (err error) {
	tx := m.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(hotspot).Error; err != nil {
		tx.Rollback()
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
	if err := tx.Create(payments).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
