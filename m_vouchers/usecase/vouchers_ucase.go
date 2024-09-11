package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type vouchersUsecase struct {
	Timeout    time.Duration
	Repository domain.VcrRepository
	Radcheck   domain.RadcheckRepository
	Package    domain.PackageRepository
}

func NewUsecase(t time.Duration, r domain.VcrRepository, radcheck domain.RadcheckRepository, pkg domain.PackageRepository) domain.VcrUsecase {
	return &vouchersUsecase{Timeout: t, Repository: r, Radcheck: radcheck, Package: pkg}
}

func (uc *vouchersUsecase) CreateBatch(ctx context.Context, pkg uint, size uint) (batch string, err error) {
	c, cancel := context.WithTimeout(ctx, uc.Timeout)
	defer cancel()

	// cari paket yang dimaksud
	pack, err := uc.Package.First(pkg)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	vouchers := []domain.Vcr{}
	radchecks := []domain.Rdcheck{}

	var timeUnix string

	for {
		// cari kode batch yang unik
		timeUnix = strconv.FormatInt(time.Now().Unix(), 10)
		vcrs, err := uc.Repository.GetByBatchname(timeUnix)
		// if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 	logrus.Error(err)
		// 	return 0,
		// }
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		if len(vcrs) != 0 {
			continue
		}

		for {
			// cari kode voucher yang unik
			voucher, err := password.Generate(6, 2, 0, true, true)
			if err != nil {
				logrus.Error(err)
				return "", err
			}
			// REPOSITORY
			vcr, err := uc.Repository.GetByUsername(c, voucher)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Error(err)
				return "", err
			}

			if vcr != (domain.Vcr{}) {
				continue
			}
			// RADCHECK
			rcs, err := uc.Radcheck.GetByUsername(c, voucher)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Error(err)
				return "", err
			}
			if len(rcs) != 0 {
				continue
			}

			if len(vouchers) < int(size) {
				vouchers = append(vouchers, domain.Vcr{Username: voucher, Pkg: pkg, Batchcode: timeUnix})
				radchecks = append(radchecks, domain.Rdcheck{Username: voucher, Attribute: "Cleartext-Password", OP: ":=", Value: voucher})
				radchecks = append(radchecks, domain.Rdcheck{Username: voucher, Attribute: "User-Profile", OP: ":=", Value: pack.Profile})
			} else {
				err = uc.Repository.Batch(c, vouchers, radchecks)
				if err != nil {
					logrus.Error(err)
					return "", err
				}
				return timeUnix, err
			}
		}

	}
}
