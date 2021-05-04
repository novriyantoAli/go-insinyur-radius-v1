package usecase

import (
	"context"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type schedulerUsecase struct {
	Timeout            time.Duration
	RepositoryRadcheck domain.RadcheckRepository
	RepositoryRadacct  domain.RadacctRepository
}

// NewSchedulerUsecase ...
func NewUsecase(t time.Duration, r domain.RadcheckRepository, a domain.RadacctRepository) domain.SchedulerUsecase {
	return &schedulerUsecase{Timeout: t, RepositoryRadcheck: r, RepositoryRadacct: a}
}

// Fetch ...
func (uc *schedulerUsecase) GetUsers(c context.Context, delete bool) (resArr []domain.Radcheck, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	resArr, err = uc.RepositoryRadcheck.FetchWithValueExpiration(ctx, delete)
	if err != nil {
		logrus.Error(err)
	}

	return

}

func (uc *schedulerUsecase) DeleteExpireUsers(c context.Context, username string) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.RepositoryRadcheck.DeleteWithUsername(ctx, username)

	return
}

func (uc *schedulerUsecase) GetOnlineUsers(c context.Context, usernameList string) (res []domain.Radacct, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.RepositoryRadacct.FetchWithUsernameBatch(ctx, usernameList)

	return
}
