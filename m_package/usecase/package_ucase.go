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

func (uc *packageUsecase) Find() (res []domain.Pkg, err error) {

	pkg := domain.Pkg{}
	res, err = uc.Repository.Find(&pkg)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *packageUsecase) First(id uint) (res domain.Pkg, err error) {
	res, err = uc.Repository.First(id)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *packageUsecase) Save(pkg *domain.Pkg) (err error) {
	err = uc.Repository.Save(pkg)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *packageUsecase) Create(c context.Context, pkg *domain.Pkg) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Create(ctx, pkg)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *packageUsecase) Delete(id uint) (err error) {
	pkgs := []domain.Pkg{{ID: id}}

	err = uc.Repository.Delete(pkgs)
	if err != nil {
		logrus.Error(err)
	}

	return
}

// func (u *packageUsecase) Fetch(c context.Context, id int64, limit int64) (res domain.PackagePage, err error) {
// 	ctx, cancel := context.WithTimeout(c, u.Timeout)
// 	defer cancel()

// 	pac := domain.Package{}
// 	totalPage, err := u.Repository.CountPage(ctx, pac)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}

// 	data, err := u.Repository.Fetch(ctx, id, limit)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}

// 	res.TotalPage = totalPage
// 	res.Data = data

// 	return
// }
