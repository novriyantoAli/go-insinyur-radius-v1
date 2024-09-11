package usecase

import (
	"context"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type customerUsecase struct {
	Timeout    time.Duration
	Repository domain.CustomerRepository
}

func NewUsecase(timeout time.Duration, r domain.CustomerRepository) domain.CustomerUsecase {
	return &customerUsecase{Timeout: timeout, Repository: r}
}

func (uc *customerUsecase) Find(c context.Context, limit uint) (res []domain.Customer, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	customer := domain.Customer{}
	res, err = uc.Repository.Find(ctx, &customer)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *customerUsecase) Save(c context.Context, cr *domain.CustomerRequest) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	format := "1/2/2006"
	birthdate, err := time.Parse(format, cr.Birthdate)
	if err != nil {
		logrus.Error(err)
		return err
	}

	cst := domain.Customer{
		Name:      cr.Name,
		Identity:  cr.Identity,
		Sex:       cr.Sex,
		Email:     cr.Email,
		Phone:     cr.Phone,
		Job:       cr.Job,
		Image:     cr.Image,
		Birthdate: birthdate,
		Lat:       cr.Lat,
		Lng:       cr.Lng,
	}

	err = uc.Repository.Save(ctx, &cst)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *customerUsecase) Delete(request *domain.CustomerRequestUpDel) (err error) {
	err = uc.Repository.Delete(&domain.Customer{
		ID:       request.ID,
		Name:     request.Name,
		Identity: request.Identity,
		Sex:      request.Sex,
		Email:    request.Email,
		Phone:    request.Phone,
		Job:      request.Job,
		Image:    request.Image,
	})
	if err != nil {
		logrus.Error(err)
	}

	return
}
