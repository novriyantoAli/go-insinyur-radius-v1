package mysql

import (
	"context"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type mysqlRepository struct {
	Conn *gorm.DB
}

// NewMysqlRepository ...
func NewMysqlRepository(conn *gorm.DB) domain.RadreplyRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) First(ctx context.Context, radreply *domain.Radreply) (res domain.Radreply, err error) {
	// err = m.Conng.Find(&pkgs, pkg).Error
	err = m.Conn.First(&res, radreply).Error
	// err = m.Conng.Model(pkg).Preload("Radusergroup").Find(&pkgs).Error
	return
}

func (m *mysqlRepository) Delete(ctx context.Context, radreply *domain.Radreply) (err error) {
	err = m.Conn.Delete(&domain.Radreply{}, "username LIKE ?", "%"+radreply.Username+"%").Error
	return
}
