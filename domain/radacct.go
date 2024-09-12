package domain

import (
	"context"
	"time"
)

type Rdacct struct {
	Radacctid          uint      `json:"radacctid"`
	Acctsessionid      string    `json:"acctsessionid"`
	Acctuniqueid       string    `json:"acctuniqueid"`
	Username           string    `json:"username"`
	Realm              string    `json:"realm"`
	Nasipaddress       string    `json:"nasipaddress"`
	Nasportid          string    `json:"nasportid"`
	Nasporttype        string    `json:"nasporttype"`
	Acctstarttime      time.Time `json:"acctstarttime"`
	Acctupdatetime     time.Time `json:"acctupdatetime"`
	Acctstoptime       time.Time `json:"acctstoptime"`
	Acctinterval       int64     `json:"acctinterval"`
	Acctsessiontime    int64     `json:"acctsessiontime"`
	Acctauthentic      string    `json:"acctauthentic"`
	ConnectinfoStart   string    `json:"connectinfo_start"`
	ConnectinfoStop    string    `json:"connectinfo_stop"`
	Acctinputoctets    int64     `json:"acctinputoctets"`
	Acctoutputoctets   int64     `json:"acctoutputoctets"`
	Calledstationid    string    `json:"calledstationid"`
	Callingstationid   string    `json:"callingstationid"`
	Acctterminatecause string    `json:"acctterminatecause"`
	Servicetype        string    `json:"servicetype"`
	Framedprotocol     string    `json:"framedprotocol"`
	Framedipaddress    string    `json:"framedipaddress"`
	Nas                Nas       `json:"nas,omitempty" gorm:"foreignKey:Nasipaddress;references:Nasname"`
}

func (Rdacct) TableName() string {
	return "radacct"
}

// Radacct ...
type Radacct struct {
	Radacctid          *int64     `json:"radacctid"`
	Acctsessionid      *string    `json:"acctsessionid"`
	Acctuniqueid       *string    `json:"acctuniqueid"`
	Username           *string    `json:"username"`
	Realm              *string    `json:"realm"`
	Nasipaddress       *string    `json:"nasipaddress"`
	Nasportid          *string    `json:"nasportid"`
	Nasporttype        *string    `json:"nasporttype"`
	Acctstarttime      *time.Time `json:"acctstarttime"`
	Acctupdatetime     *time.Time `json:"acctupdatetime"`
	Acctstoptime       *time.Time `json:"acctstoptime"`
	Acctinterval       *int64     `json:"acctinterval"`
	Acctsessiontime    *int64     `json:"acctsessiontime"`
	Acctauthentic      *string    `json:"acctauthentic"`
	ConnectinfoStart   *string    `json:"connectinfo_start"`
	ConnectinfoStop    *string    `json:"connectinfo_stop"`
	Acctinputoctets    *int64     `json:"acctinputoctets"`
	Acctoutputoctets   *int64     `json:"acctoutputoctets"`
	Calledstationid    *string    `json:"calledstationid"`
	Callingstationid   *string    `json:"callingstationid"`
	Acctterminatecause *string    `json:"acctterminatecause"`
	Servicetype        *string    `json:"servicetype"`
	Framedprotocol     *string    `json:"framedprotocol"`
	Framedipaddress    *string    `json:"framedipaddress"`
	Secret             *string    `json:"secret"`
}

// RadacctUsecase ...
type RadacctUsecase interface {
	FetchWithUsernameBatch(ctx context.Context, usernameList string) (res []Radacct, err error)
}

// RadacctRepository ...
type RadacctRepository interface {
	FetchWithUsernameBatch(ctx context.Context, usernameList string) (res []Radacct, err error)
	Get(ctx context.Context, radacct Radacct) (res []Radacct, err error)
	Find(ctx context.Context, param *Rdacct) (res []Rdacct, err error)
	Save(ctx context.Context, param *Rdacct) (err error)
	FindUsernameIn(ctx context.Context, usernamein []string) (res []Rdacct, err error)
	FindToday(ctx context.Context, todaystart string, todayend string) (res []Rdacct, err error)
	FindWeek(ctx context.Context, weekbreakdown int) (res []Rdacct, err error)
	FindMonth(ctx context.Context, targetmonth int, targetyear int) (res []Rdacct, err error)
}
