package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type paymentUsecase struct {
	Timeout    time.Duration
	Repository domain.PaymentRepository
}

func NewUsecase(timeout time.Duration, r domain.PaymentRepository) domain.PaymentUsecase {
	return &paymentUsecase{Timeout: timeout, Repository: r}
}

func (uc *paymentUsecase) CurrentYear(c context.Context) (res []domain.DashboardYearResp, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	last12month := time.Now().AddDate(0, -11, 0)
	for i := 0; i < 12; i++ {
		tmx := last12month.AddDate(0, i, 0)
		resArr, err := uc.Repository.Month(ctx, int(tmx.Month()), tmx.Year())
		if err != nil {
			return res, err
		}
		// earning == credit else dbt
		lcl := domain.DashboardYearResp{
			Month:    tmx.Format("Jan 2006"),
			Earning:  0.0,
			Expenses: 0.0,
		}
		for _, v := range resArr {
			if v.Status == "kredit" {
				lcl.Earning += v.Amount
			} else {
				lcl.Expenses += v.Amount
			}
		}
		res = append(res, lcl)
	}

	return
}

func (uc *paymentUsecase) CurrentMonth(c context.Context) (res []domain.Payment, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Month(ctx, int(time.Now().Month()), time.Now().Year())
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *paymentUsecase) TodayReport(c context.Context) (res []domain.Payment, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	now := time.Now()

	todayStart := now.Format("2006-01-02")
	todayStart += " 00:00:00"

	todayEnd := now.Format("2006-01-02")
	todayEnd += " 23:59:59"

	res, err = uc.Repository.Today(ctx, todayStart, todayEnd)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *paymentUsecase) PaymentOrder(c context.Context, idbatch string, amount uint64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	payment := domain.Payment{IDBatch: idbatch}

	res, err := uc.Repository.Find(ctx, &payment)
	if err != nil {
		logrus.Error(err)
	}

	if len(res) == 0 {
		return errors.New("order id not found")
	} else {
		kredit := 0
		debit := 0
		for _, value := range res {
			if value.Status == "kredit" {
				kredit += int(value.Amount)
			} else {
				debit += int(value.Amount)
			}
		}

		if (kredit - debit) == 0 {
			return errors.New("order sudah lunas")
		}

		if ((kredit - debit) - int(amount)) == 0 {
			payment.Amount = float64(amount)
			payment.Status = "debit"
			err = uc.Repository.PaymentOrder(ctx, &payment)
			return
		}

		if ((kredit - debit) - int(amount)) < 0 {
			return errors.New("jumlah yang dibayarkan sudah berlebih")
		}

		return
	}
}
