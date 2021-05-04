package usecase

import (
	"context"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type packageUsecase struct {
	Timeout    time.Duration
	Repository domain.PackageRepository
}

func NewUsecase(timeout time.Duration, r domain.PackageRepository) domain.PackageUsecase {
	return &packageUsecase{Timeout: timeout, Repository: r}
}

func (u *packageUsecase) Fetch(c context.Context, id int64, limit int64) (res domain.PackagePage, err error) {
	ctx, cancel := context.WithTimeout(c, u.Timeout)
	defer cancel()

	pac := domain.Package{}
	totalPage, err := u.Repository.CountPage(ctx, pac)
	if err != nil {
		logrus.Error(err)
		return
	}

	data, err := u.Repository.Fetch(ctx, id, limit)
	if err != nil {
		logrus.Error(err)
		return
	}

	res.TotalPage = totalPage
	res.Data = data

	return
}
