package domain

import "context"

type Radcheck struct {
	ID        *int64  `json:"id"`
	Username  *string `json:"username"`
	Attribute *string `json:"attribute"`
	OP        *string `json:"op"`
	Value     *string `json:"value"`
}

type Rdcheck struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Attribute string
	OP        string
	Value     string
}

type Tabler interface {
	TableName() string
}

func (Rdcheck) TableName() string {
	return "radcheck"
}

type RadcheckRepository interface {
	First(ctx context.Context, radcheck *Rdcheck) (res Rdcheck, err error)
	FindExpireToday(ctx context.Context) (res []Rdcheck, err error)
	FindExpireTwoWeek(ctx context.Context) (res []Rdcheck, err error)
	Delete(ctx context.Context, radcheck *Rdcheck) (err error)
	GetByUsername(ctx context.Context, username string) (radchecks []Rdcheck, err error)
	Get(ctx context.Context, radcheck Radcheck) (res []Radcheck, err error)
	FetchWithValueExpiration(ctx context.Context, delete bool) (res []Radcheck, err error)
	Update(ctx context.Context, radcheck Radcheck) (err error)
	DeleteWithUsername(ctx context.Context, username string) (err error)
}
