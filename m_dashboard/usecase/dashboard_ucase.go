package usecase

import (
	"context"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/helper"
	"github.com/sirupsen/logrus"
)

type dashboardUsecase struct {
	Timeout time.Duration
	CSR     domain.CustomerRepository
	CLR     domain.ClientRepository
	USR     domain.UsersRepository
	PKR     domain.PackageRepository
	TLR     domain.TimelineRepository
	RAR     domain.RadacctRepository
}

func NewUsecase(tout time.Duration, csr domain.CustomerRepository, clr domain.ClientRepository, usr domain.UsersRepository, pkr domain.PackageRepository, tlr domain.TimelineRepository, rar domain.RadacctRepository) domain.DashboardUsecase {
	return &dashboardUsecase{Timeout: tout, CSR: csr, CLR: clr, USR: usr, PKR: pkr, TLR: tlr, RAR: rar}
}

func (uc *dashboardUsecase) TopLocation(c context.Context) (res []domain.DashboardTopLocation, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	now := time.Now()

	todayStart := now.Format("2006-01-02")
	todayStart += " 00:00:00"

	todayEnd := now.Format("2006-01-02")
	todayEnd += " 23:59:59"

	resToday, err := uc.RAR.FindToday(ctx, todayStart, todayEnd)
	if err != nil {
		logrus.Error(err)
		return
	}

	groupByCalledstationid := helper.GroupByProperty(resToday, func(radacct domain.Rdacct) string {
		return radacct.Calledstationid //o.OrderHotspot.IDPackage
	})

	for calledstationid, group := range groupByCalledstationid {
		res = append(res, domain.DashboardTopLocation{
			Name:          calledstationid,
			CountingToday: len(group),
		})
	}

	resWeek, err := uc.RAR.FindWeek(ctx, 1)
	if err != nil {
		logrus.Error(err)
		return
	}

	groupByCalledstationid2 := helper.GroupByProperty(resWeek, func(radacct domain.Rdacct) string {
		return radacct.Calledstationid
	})

	for calledstationid, group := range groupByCalledstationid2 {
		for i := 0; i < len(res); i++ {
			if calledstationid == res[i].Name {
				res[i].CountingWeek = len(group)
			}
		}
	}

	resMonth, err := uc.RAR.FindMonth(ctx, int(now.Month()), now.Year())
	if err != nil {
		logrus.Error(err)
		return
	}

	groupByCalledstationid3 := helper.GroupByProperty(resMonth, func(radacct domain.Rdacct) string {
		return radacct.Calledstationid
	})

	for calledstationid, group := range groupByCalledstationid3 {
		for i := 0; i < len(res); i++ {
			if calledstationid == res[i].Name {
				res[i].CountingMonth = len(group)
			}
		}
	}

	return

}

func (uc *dashboardUsecase) Timeline(c context.Context) (res []domain.Timeline, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.TLR.Find(ctx)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *dashboardUsecase) QuickCount(c context.Context) (res []domain.DashboardQuick, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	customers, err := uc.CSR.Find(ctx, &domain.Customer{})
	if err != nil {
		logrus.Error(err)
		return
	}

	customer := domain.DashboardQuick{
		ID:              "pelanggan",
		Title:           "Pelanggan",
		Value:           len(customers),
		Icon:            "mso-attach_money",
		ChangeText:      "0",
		ChangeDirection: "down",
	}

	res = append(res, customer)

	clients, err := uc.CLR.Find(&domain.Client{})
	if err != nil {
		logrus.Error(err)
		return
	}

	client := domain.DashboardQuick{
		ID:              "klien",
		Title:           "Klien",
		Value:           len(clients),
		Icon:            "mso-attach_money",
		ChangeText:      "0",
		ChangeDirection: "down",
	}

	res = append(res, client)

	users, err := uc.USR.Find(&domain.Usr{})
	if err != nil {
		logrus.Error(err)
		return
	}

	user := domain.DashboardQuick{
		ID:              "pengguna",
		Title:           "Pengguna",
		Value:           len(users),
		Icon:            "mso-attach_money",
		ChangeText:      "0",
		ChangeDirection: "down",
	}

	res = append(res, user)

	packages, err := uc.PKR.Find(&domain.Pkg{})
	if err != nil {
		logrus.Error(err)
		return
	}

	pkg := domain.DashboardQuick{
		ID:              "package_hotspot",
		Title:           "Paket Hotspot",
		Value:           len(packages),
		Icon:            "mso-attach_money",
		ChangeText:      "0",
		ChangeDirection: "down",
	}

	res = append(res, pkg)

	return
}
