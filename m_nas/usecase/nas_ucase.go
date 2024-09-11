package usecase

import (
	"context"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type nasUsecase struct {
	Repository domain.NasRepository
	Timeout    time.Duration
}

func NewUsecase(repository domain.NasRepository, timeout time.Duration) domain.NasUsecase {
	return &nasUsecase{Repository: repository, Timeout: timeout}
}

func (uc *nasUsecase) Find(c context.Context, nas *domain.Nas) (res []domain.Nas, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Find(ctx, nas)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *nasUsecase) Create(c context.Context, req *domain.NasReq) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	param := domain.Nas{
		Nasname:   req.Nasname,
		Shortname: req.Shortname,
		Type:      req.Type,
		Secret:    req.Secret,
	}

	err = uc.Repository.Create(ctx, &param)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *nasUsecase) Delete(c context.Context, req *domain.NasDeleteReq) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	param := domain.Nas{
		ID:        req.ID,
		Nasname:   req.Nasname,
		Shortname: req.Shortname,
		Type:      req.Type,
		Secret:    req.Secret,
	}

	err = uc.Repository.Delete(ctx, &param)
	if err != nil {
		logrus.Error(err)
	}

	return
}
