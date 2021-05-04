package usecase

import (
	"context"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type transactionUsecase struct {
	Timeout    time.Duration
	Repository domain.TransactionRepository
}

func NewUsecase(timeout time.Duration, r domain.TransactionRepository) domain.TransactionUsecase {
	return &transactionUsecase{Timeout: timeout, Repository: r}
}

func (u *transactionUsecase) Fetch(c context.Context, idUsers int64, id int64, limit int64, query string) (res domain.TransactionPage, err error) {
	ctx, cancel := context.WithTimeout(c, u.Timeout)
	defer cancel()

	transaction := domain.Transaction{}
	transaction.IDReseller = &idUsers
	totalPage, err := u.Repository.CountPage(ctx, transaction)
	if err != nil {
		logrus.Error(err)
		return
	}

	data, err := u.Repository.Fetch(ctx, idUsers, id, limit)
	if err != nil {
		logrus.Error(err)
		return
	}

	res.TotalPage = totalPage
	res.Data = data

	return
}
