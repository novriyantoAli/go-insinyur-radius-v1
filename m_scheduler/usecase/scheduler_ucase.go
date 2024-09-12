package usecase

import (
	"bytes"
	"context"
	"math"
	"os/exec"
	"strings"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type schedulerUsecase struct {
	Timeout            time.Duration
	RepositoryRadcheck domain.RadcheckRepository
	RepositoryRadacct  domain.RadacctRepository
	ROrder             domain.OrderRepository
	RIPBinding         domain.IPBindingRepository
	RFirewall          domain.FirewallRepository
	RSimpleQueue       domain.SimpleQueueRepository
	ROSClient          []domain.ROSClient
}

// NewSchedulerUsecase ...
func NewUsecase(t time.Duration, r domain.RadcheckRepository, a domain.RadacctRepository, rorder domain.OrderRepository, ripbinding domain.IPBindingRepository, rfirewall domain.FirewallRepository, rsimplequeue domain.SimpleQueueRepository, rosclient []domain.ROSClient) domain.SchedulerUsecase {
	return &schedulerUsecase{Timeout: t, RepositoryRadcheck: r, RepositoryRadacct: a, ROrder: rorder, RIPBinding: ripbinding, RFirewall: rfirewall, RSimpleQueue: rsimplequeue, ROSClient: rosclient}
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

func (uc *schedulerUsecase) IncreasePerform(c context.Context, wita *time.Location) (err error) {
	// find 15 day expire
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	layoutFormat := "02 Jan 2006 15:04:05"

	res, err := uc.RepositoryRadcheck.FindExpireTwoWeek(ctx)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if len(res) <= 0 {
		logrus.Infoln("two week increase Perform clear")
		return
	}

	targetUsername := []string{}
	for _, radcheck := range res {
		targetUsername = append(targetUsername, radcheck.Username)
	}

	resacct, err := uc.RepositoryRadacct.FindUsernameIn(ctx, targetUsername)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, radcheck := range res {
		date, err := time.ParseInLocation(layoutFormat, radcheck.Value, wita)
		if err != nil {
			logrus.Error(err)
		}

		duration := time.Now().Sub(date)
		if math.Signbit(duration.Seconds()) == false {
			for _, radacct := range resacct {
				if radcheck.Username == radacct.Username {
					command := `echo "Acct-Session-Id=` + radacct.Acctsessionid + `,User-Name=` + radacct.Username + `,NAS-IP-Address=` + radacct.Nasipaddress + `,Framed-IP-Address=` + radacct.Framedipaddress + `" | radclient -x ` + radacct.Nasipaddress + `:3799 disconnect '` + radacct.Nas.Secret + `'`
					cmd := exec.Command("sh", "-c", command)
					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr
					cmd.Run()
					// set date to radacct
					radacct.Acctstoptime = time.Now()
					err = uc.RepositoryRadacct.Save(ctx, &radacct)
					if err != nil {
						logrus.Error(err)
						err = nil
					}
				}
			}
			// delete expire user
			err = uc.RepositoryRadcheck.Delete(ctx, &radcheck)
			if err != nil {
				logrus.Error(err)
			}
		}
	}

	return
}

func (uc *schedulerUsecase) GuardRadius(c context.Context, wita *time.Location) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	layoutFormat := "02 Jan 2006 15:04:05"

	res, err := uc.RepositoryRadcheck.FindExpireToday(ctx)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if len(res) <= 0 {
		logrus.Infoln("not user expire today")
		return
	}

	targetUsername := []string{}
	for _, radcheck := range res {
		targetUsername = append(targetUsername, radcheck.Username)
	}

	resacct, err := uc.RepositoryRadacct.FindUsernameIn(ctx, targetUsername)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, radcheck := range res {
		date, err := time.ParseInLocation(layoutFormat, radcheck.Value, wita)
		if err != nil {
			logrus.Error(err)
		}

		duration := time.Now().Sub(date)
		if math.Signbit(duration.Seconds()) == false {
			for _, radacct := range resacct {
				if radcheck.Username == radacct.Username {
					command := `echo "Acct-Session-Id=` + radacct.Acctsessionid + `,User-Name=` + radacct.Username + `,NAS-IP-Address=` + radacct.Nasipaddress + `,Framed-IP-Address=` + radacct.Framedipaddress + `" | radclient -x ` + radacct.Nasipaddress + `:3799 disconnect '` + radacct.Nas.Secret + `'`
					cmd := exec.Command("sh", "-c", command)
					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr
					cmd.Run()
					radacct.Acctstoptime = time.Now()
					err = uc.RepositoryRadacct.Save(ctx, &radacct)
					if err != nil {
						logrus.Error(err)
						err = nil
					}
				}
			}

		}
	}

	return
}

func (uc *schedulerUsecase) GuardMacBinding(c context.Context, wita *time.Location) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	layoutFormat := "02 Jan 2006 15:04:05"

	res, err := uc.ROrder.FindOrderClientExpire(ctx)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, orderClient := range res {
		date, _ := time.ParseInLocation(layoutFormat, orderClient.Expire, wita)
		duration := time.Now().Sub(date)
		if math.Signbit(duration.Seconds()) == false {
			// search in database
			for _, clientBinding := range orderClient.Client.ClientBinding {
				macAddress := strings.ToUpper(clientBinding.Mac)
				// delete ip binding
				for _, value := range uc.ROSClient {
					ipbinding, err := uc.RIPBinding.Firsts(value.Client, &domain.IPBinding{MacAddress: macAddress})
					if err != nil {
						logrus.Error(err)
					}
					err = uc.RIPBinding.Deletes(value.Client, &ipbinding)
					if err != nil {
						logrus.Error(err)
					}

					firewallConn, err := uc.RFirewall.Firsts(value.Client, &domain.Firewall{SrcMacAddress: macAddress})
					if err != nil {
						logrus.Error(err)
					}

					firewallPack, err := uc.RFirewall.Firsts(value.Client, &domain.Firewall{ConnectionMark: firewallConn.NewConnectionMark})
					if err != nil {
						logrus.Error(err)
					}

					// delete queue
					simplequeue, err := uc.RSimpleQueue.Firsts(value.Client, &domain.SimpleQueue{Name: firewallPack.NewPacketMark})
					if err != nil {
						logrus.Error(err)
					}

					err = uc.RFirewall.Deletes(value.Client, &firewallConn)
					if err != nil {
						logrus.Error(err)
					}

					err = uc.RFirewall.Deletes(value.Client, &firewallPack)
					if err != nil {
						logrus.Error(err)
					}

					err = uc.RSimpleQueue.Deletes(value.Client, &simplequeue)
					if err != nil {
						logrus.Error(err)
					}

				}

				// delete firewall
				// firewallConn, err := uc.RFirewall.First(&domain.Firewall{SrcMacAddress: macAddress})
				// if err != nil {
				// 	logrus.Error(err)
				// }

				// firewallPack, err := uc.RFirewall.First(&domain.Firewall{ConnectionMark: firewallConn.NewConnectionMark})
				// if err != nil {
				// 	logrus.Error(err)
				// }

				// // delete queue
				// simplequeue, err := uc.RSimpleQueue.First(&domain.SimpleQueue{Name: firewallPack.NewPacketMark})
				// if err != nil {
				// 	logrus.Error(err)
				// }

				// err = uc.RFirewall.Delete(&firewallConn)
				// if err != nil {
				// 	logrus.Error(err)
				// }

				// err = uc.RFirewall.Delete(&firewallPack)
				// if err != nil {
				// 	logrus.Error(err)
				// }

				// err = uc.RSimpleQueue.Delete(&simplequeue)
				// if err != nil {
				// 	logrus.Error(err)
				// }
			}
		}
	}

	return

}
