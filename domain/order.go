package domain

import (
	"context"
	"time"
)

type Order struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	IDUser       uint         `json:"id_user"`
	IDCustomer   uint         `json:"id_customer"`
	IDBatch      string       `json:"id_batch" gorm:"index;unique"`
	Amount       uint         `json:"amount"`
	User         Usr          `json:"user,omitempty" gorm:"foreignKey:IDUser;references:ID"`
	Customer     Customer     `json:"customer,omitempty" gorm:"foreignKey:IDCustomer;references:ID"`
	OrderHotspot OrderHotspot `json:"order_hotspot,omitempty" gorm:"foreignKey:IDBatch;references:IDBatch"`
	OrderClient  OrderClient  `json:"order_client,omitempty" gorm:"foreignKey:IDBatch;references:IDBatch"`
	Payment      []Payment    `json:"payment" json:"vouchers" gorm:"foreignKey:IDBatch;references:IDBatch"`
	CreatedAt    time.Time    `json:"created_at"`
}

type OrderHotspot struct {
	IDBatch                    string `gorm:"index;unique" json:"id_batch"`
	IDPackage                  uint   `json:"id_package"`
	IncludeSymbols             bool   `json:"include_symbols"`
	IncludeNumbers             bool   `json:"include_numbers"`
	IncludeLowercaseLetters    bool   `json:"include_lowercase_letters"`
	IncludeUppercaseLetters    bool   `json:"include_uppercase_letters"`
	ExcludeSimilarCharacters   bool   `json:"exclude_similar_characters"`
	ExcludeAmbiguousCharacters bool   `json:"exclude_ambiguous_characters"`
	Activated                  bool   `json:"activated"`
	Package                    Pkg    `json:"package,omitempty" gorm:"foreignKey:IDPackage;references:ID"`
	Vouchers                   []Vcr  `json:"vouchers" gorm:"foreignKey:Batchcode;references:IDBatch"`
}

type OrderClient struct {
	IDBatch   string `gorm:"index;unique" json:"id_batch"`
	IDClient  uint   `gorm:"index" json:"id_client"`
	Client    Client `json:"client,omitempty" gorm:"foreignKey:IDClient;references:ID"`
	Expire    string `json:"expire"`
	OrderDate time.Time
}

type UserOrder struct {
	IDBatch    string `json:"id_batch" gorm:"index"`
	IDUser     uint   `json:"id_user"`
	IDCustomer uint   `json:"id_customer"`
}

type OrderClientRequest struct {
	Payment  string `json:"payment" validate:"required"`
	IDUser   uint   `json:"id_user" validate:"required"`
	IDClient uint   `json:"id_client" validate:"required"`
	Amount   string `json:"amount" validate:"required"`
}

type OrderReportResp struct {
	Hotspot int `json:"hotspot"`
	Client  int `json:"client"`
}

type OrderHotspotMonthReport struct {
	MonthBreakdown int    `json:"month_breakdown"`
	MonthCurrent   int    `json:"month_current"`
	Winning        string `json:"winning"`
	WinningData    int    `json:"winning_data"`
	OtherData      int    `json:"other_data"`
}

type OrderClientMonthReport struct {
	MonthBreakdown   int      `json:"month_breakdown"`
	MonthCurrent     int      `json:"month_current"`
	MonthName        []string `json:"month_name"`
	MonthAchievement []int    `json:"month_achievement"`
}

type OrderRequest struct {
	IDUser                     uint   `json:"id_user" validate:"required"`
	IDCustomer                 uint   `json:"id_customer" validate:"required"`
	IDPackage                  uint   `json:"id_package" validate:"required"`
	IncludeSymbols             bool   `json:"include_symbols"`
	IncludeNumbers             bool   `json:"include_numbers"`
	IncludeLowercaseLetters    bool   `json:"include_lowercase_letters"`
	IncludeUppercaseLetters    bool   `json:"include_uppercase_letters"`
	ExcludeSimilarCharacters   bool   `json:"exclude_similar_characters"`
	ExcludeAmbiguousCharacters bool   `json:"exclude_ambiguous_characters"`
	Activated                  bool   `json:"activated"`
	Amount                     string `json:"amount" validate:"required"`
	Payment                    string `json:"payment" validate:"required"`
}

type OrderRepository interface {
	Find(ctx context.Context, order *Order) (res []Order, err error)
	FindMonthCreated(ctx context.Context, targetmonth int, targetyear int) (res []Order, err error)
	FindOrderClientExpire(ctx context.Context) (res []OrderClient, err error)
	FindOrderClient(ctx context.Context, oc *OrderClient) (res []OrderClient, err error)
	Paginate(ctx context.Context, lastID uint, limit int) (res []Order, err error)
	First(ctx context.Context, order *Order) (res Order, err error)
	Client(ctx context.Context, order *Order, orderClient *OrderClient, radcheck []Rdcheck, radreply []Radreply, payments []Payment) (err error)
	Batch(ctx context.Context, order *Order, hotspot *OrderHotspot, vcrs []Vcr, radcheck []Rdcheck, payments []Payment) (err error)
}

type OrderUsecase interface {
	Find(ctx context.Context, lastID uint, limit int) (res []Order, err error)
	ReportHotspotMonth(ctx context.Context) (res OrderHotspotMonthReport, err error)
	ReportClientMonth(ctx context.Context) (ers OrderClientMonthReport, err error)
	Client(ctx context.Context, request *OrderClientRequest) (err error)
	Order(ctx context.Context, request *OrderRequest) (err error)
}
