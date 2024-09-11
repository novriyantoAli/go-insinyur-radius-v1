package domain

import (
	"context"
	// "gorm.io/gorm"
)

type Radusergrouprequest struct {
	Username      string          `json:"username" validate:"required"`
	Groupname     string          `json:"groupname" validate:"required"`
	Priority      uint            `json:"priority" validate:"required"`
	Radgroupcheck []Radgroupcheck `json:"Radgroupcheck"`
	Radgroupreply []Radgroupreply `json:"Radgroupreply"`
}

type Radusergroup struct {
	// gorm.Model
	Username      string          `json:"username" gorm:"index;unique"`
	Groupname     string          `json:"groupname"`
	Priority      uint            `json:"priority"`
	Radgroupcheck []Radgroupcheck `gorm:"foreignKey:Groupname;references:Groupname"`
	Radgroupreply []Radgroupreply `gorm:"foreignKey:Groupname;references:Groupname"`
}

func (Radusergroup) TableName() string {
	return "radusergroup"
}

type ProfilesUsecase interface {
	Fetch(c context.Context) (rug []Radusergroup, err error)
	Save(r *Radusergrouprequest) (err error)
	Delete(r *Radusergrouprequest) (err error)
	// Find(c context.Context, groupname string) (rug Radusergroup, err error)
}

type ProfilesRepository interface {
	First(r *Radusergroup) (rug Radusergroup, err error)
	Find(r *Radusergroup) (rug []Radusergroup, err error)
	Save(r *Radusergroup) (err error)
	Delete(r *Radusergroup) (err error)
}
