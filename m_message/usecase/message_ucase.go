package usecase

import (
	"context"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
)

type messageUsecase struct {
	Timeout time.Duration
	VCR     domain.VcrRepository
}

func NewUsecase(tout time.Duration, vcr domain.VcrRepository) domain.MessageUsecase {
	return &messageUsecase{Timeout: tout, VCR: vcr}
}

func (uc *messageUsecase) ReadUserMessage(c context.Context, usr *domain.User) (err error) {
	_, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	// // check if user is a voucher
	// vcr, err := uc.VCR.First(ctx, &domain.Vcr{Username: usr.Username})
	// _, err = nil {
	// 	logrus.Error(err)

	// 	return err
	// }

	// count if vcr

	return

}
