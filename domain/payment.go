package domain

import (
	"context"
	"time"
)

type Payment struct {
	IDBatch   string  `json:"id_batch" gorm:"index"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
	CreatedAt time.Time
}

type PayReq struct {
	Amount  string `json:"amount" validate:"required"`
	IDBatch string `json:"id_batch" validate:"required"`
}

type DashboardYearResp struct {
	Month    string  `json:"month"`
	Earning  float64 `json:"earning"`
	Expenses float64 `json:"expenses"`
}

type DebCreResp struct {
	Credit float64 `json:"credit"`
	Debt   float64 `json:"debt"`
}

type PaymentRepository interface {
	Find(ctx context.Context, payment *Payment) (res []Payment, err error)
	PaymentOrder(ctx context.Context, payment *Payment) (err error)
	Today(ctx context.Context, todaystart string, todayend string) (res []Payment, err error)
	Month(ctx context.Context, targetmonth int, targetyear int) (res []Payment, err error)
	Year(ctx context.Context, targetyear int) (res []Payment, err error)
}

type PaymentUsecase interface {
	PaymentOrder(c context.Context, idbatch string, amount uint64) (err error)
	TodayReport(ctx context.Context) (res []Payment, err error)
	CurrentMonth(c context.Context) (res []Payment, err error)
	CurrentYear(c context.Context) (res []DashboardYearResp, err error)
}
