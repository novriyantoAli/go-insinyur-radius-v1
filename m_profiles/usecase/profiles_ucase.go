package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type profileUsecase struct {
	Timeout    time.Duration
	Repository domain.ProfilesRepository
}

func NewUsecase(timeout time.Duration, r domain.ProfilesRepository) domain.ProfilesUsecase {
	return &profileUsecase{Timeout: timeout, Repository: r}
}

func (uc *profileUsecase) Fetch(ctx context.Context) (rug []domain.Radusergroup, err error) {
	_, cancel := context.WithTimeout(ctx, uc.Timeout)
	defer cancel()

	radusergroup := domain.Radusergroup{}
	rug, err = uc.Repository.Find(&radusergroup)

	return
}

func (uc *profileUsecase) Save(r *domain.Radusergrouprequest) (err error) {
	// cchange white space with
	ugname := strings.ReplaceAll(r.Username, " ", "-")
	rug := domain.Radusergroup{
		Username:      ugname,
		Groupname:     ugname,
		Priority:      r.Priority,
		Radgroupcheck: r.Radgroupcheck,
		Radgroupreply: r.Radgroupreply,
	}

	// rugs, err := uc.Repository.First(&rug)
	// if errors.Is() {

	// }
	err = uc.Repository.Save(&rug)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *profileUsecase) Delete(r *domain.Radusergrouprequest) (err error) {
	rug := domain.Radusergroup{
		Username:      r.Username,
		Groupname:     r.Groupname,
		Priority:      r.Priority,
		Radgroupcheck: r.Radgroupcheck,
		Radgroupreply: r.Radgroupreply,
	}

	err = uc.Repository.Delete(&rug)
	if err != nil {
		logrus.Error(err)
	}

	return
}
