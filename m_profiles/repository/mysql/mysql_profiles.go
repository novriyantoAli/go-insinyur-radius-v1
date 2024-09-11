package mysql

import (
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

type profileRepository struct {
	Conn *gorm.DB
}

func NewMysqlRepository(conn *gorm.DB) domain.ProfilesRepository {
	return &profileRepository{Conn: conn}
}

func (repo *profileRepository) First(r *domain.Radusergroup) (rug domain.Radusergroup, err error) {
	err = repo.Conn.Model(r).Preload("Radgroupcheck").Preload("Radgroupreply").First(&rug).Error
	return
}

func (repo *profileRepository) Find(r *domain.Radusergroup) (rug []domain.Radusergroup, err error) {
	err = repo.Conn.Model(r).Preload("Radgroupcheck").Preload("Radgroupreply").Find(&rug).Error
	return
}
func (repo *profileRepository) Save(r *domain.Radusergroup) (err error) {
	err = repo.Conn.Where("username = ?", r.Username).Save(r).Error
	return
}

func (repo *profileRepository) Delete(r *domain.Radusergroup) (err error) {
	err = repo.Conn.Where("username = ?", r.Username).Delete(&r).Error
	return
}
