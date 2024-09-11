package domain

// import "gorm.io/gorm"

type Radgroupreply struct {
	// gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	Groupname string `json:"groupname" validate:"required"`
	Attribute string `json:"attribute" validate:"required"`
	OP        string `json:"op" validate:"required"`
	Value     string `json:"value" validate:"required"`
}

type Radgroupcheck struct {
	// gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	Groupname string `json:"groupname" validate:"required"`
	Attribute string `json:"attribute" validate:"required"`
	OP        string `json:"op" validate:"required"`
	Value     string `json:"value" validate:"required"`
}

func (Radgroupreply) TableName() string {
	return "radgroupreply"
}

func (Radgroupcheck) TableName() string {
	return "radgroupcheck"
}

type RadgroupRepository interface {
	FindCheck(rgck *Radgroupcheck) (res []Radgroupcheck, err error)
	FindReply(rgry *Radgroupreply) (res []Radgroupreply, err error)
	SaveCheck(rgck *Radgroupcheck) (err error)
	SaveReply(rgry *Radgroupreply) (err error)
	DeleteCheck(rgck *Radgroupcheck) (err error)
	DeleteReply(rgry *Radgroupreply) (err error)
}

type RadgroupUsecase interface {
	FindCheck(rgck *Radgroupcheck) (res []Radgroupcheck, err error)
	FindReply(rgry *Radgroupreply) (res []Radgroupreply, err error)
	SaveCheck(rgck *Radgroupcheck) (err error)
	SaveReply(rgry *Radgroupreply) (err error)
	DeleteCheck(rgck *Radgroupcheck) (err error)
	DeleteReply(rgry *Radgroupreply) (err error)
}
