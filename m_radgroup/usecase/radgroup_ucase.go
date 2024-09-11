package usecase

import (
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type radgroupUsecase struct {
	Repository domain.RadgroupRepository
}

func NewUsecase(r domain.RadgroupRepository) domain.RadgroupUsecase {
	return &radgroupUsecase{Repository: r}
}

func (uc *radgroupUsecase) FindCheck(rgck *domain.Radgroupcheck) (res []domain.Radgroupcheck, err error) {
	res, err = uc.Repository.FindCheck(rgck)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *radgroupUsecase) FindReply(rgry *domain.Radgroupreply) (res []domain.Radgroupreply, err error) {
	res, err = uc.Repository.FindReply(rgry)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *radgroupUsecase) SaveCheck(rgck *domain.Radgroupcheck) (err error) {
	err = uc.Repository.SaveCheck(rgck)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *radgroupUsecase) SaveReply(rgry *domain.Radgroupreply) (err error) {
	err = uc.Repository.SaveReply(rgry)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *radgroupUsecase) DeleteCheck(r *domain.Radgroupcheck) (err error) {
	err = uc.Repository.DeleteCheck(r)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *radgroupUsecase) DeleteReply(r *domain.Radgroupreply) (err error) {
	err = uc.Repository.DeleteReply(r)
	if err != nil {
		logrus.Error(err)
	}

	return
}
