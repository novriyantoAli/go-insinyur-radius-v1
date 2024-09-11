package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/m1/go-generate-password/generator"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/novriyantoAli/go-insinyur-radius-v1/helper"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type orderUsecase struct {
	Timeout      time.Duration
	Repository   domain.OrderRepository
	RPackage     domain.PackageRepository
	RCustomer    domain.CustomerRepository
	RVcr         domain.VcrRepository
	RRadcheck    domain.RadcheckRepository
	RRadreply    domain.RadreplyRepository
	RClient      domain.ClientRepository
	RIPBinding   domain.IPBindingRepository
	RFirewall    domain.FirewallRepository
	RSimpleQueue domain.SimpleQueueRepository
	ROSClient    []domain.ROSClient
}

func NewUsecase(timeout time.Duration, r domain.OrderRepository, rpackage domain.PackageRepository, rcustomer domain.CustomerRepository, rvcr domain.VcrRepository, rradcheck domain.RadcheckRepository, rradreply domain.RadreplyRepository, rclient domain.ClientRepository, ripbinding domain.IPBindingRepository, rfirewall domain.FirewallRepository, rsimplequeue domain.SimpleQueueRepository, rosclient []domain.ROSClient) domain.OrderUsecase {
	return &orderUsecase{Timeout: timeout, Repository: r, RPackage: rpackage, RCustomer: rcustomer, RVcr: rvcr, RRadcheck: rradcheck, RRadreply: rradreply, RClient: rclient, RIPBinding: ripbinding, RFirewall: rfirewall, RSimpleQueue: rsimplequeue, ROSClient: rosclient}
}

func (uc *orderUsecase) ReportClientMonth(c context.Context) (res domain.OrderClientMonthReport, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	year := time.Now().Year()
	month := int(time.Now().Month())

	resArr, err := uc.Repository.FindMonthCreated(ctx, month, year)
	if err != nil {
		logrus.Error(err)
	}

	resultM := slices.DeleteFunc(resArr, func(i domain.Order) bool { return i.OrderClient.IDBatch == "" })

	res.MonthCurrent = len(resultM)

	if month == 1 {
		month = 12
		year = year - 1

		resArr2, err := uc.Repository.FindMonthCreated(ctx, month, year)
		if err != nil {
			logrus.Error(err)
		}

		resultM2 := slices.DeleteFunc(resArr2, func(i domain.Order) bool { return i.OrderClient.IDBatch == "" })

		res.MonthBreakdown = len(resultM2)
	} else {

		month = month - 1
		resArr2, err := uc.Repository.FindMonthCreated(ctx, month, year)
		if err != nil {
			logrus.Error(err)
		}

		resultM2 := slices.DeleteFunc(resArr2, func(i domain.Order) bool { return i.OrderClient.IDBatch == "" })

		res.MonthBreakdown = len(resultM2)
	}

	currentTime := time.Now()

	last12Month := currentTime.AddDate(0, -11, 0)

	for i := 0; i < 12; i++ {
		month := last12Month.AddDate(0, i, 0)
		res.MonthName = append(res.MonthName, month.Month().String())
		resArr2, err := uc.Repository.FindMonthCreated(ctx, int(month.Month()), month.Year())
		if err != nil {
			logrus.Error(err)
		}
		resultM2 := slices.DeleteFunc(resArr2, func(i domain.Order) bool { return i.OrderClient.IDBatch == "" })
		res.MonthAchievement = append(res.MonthAchievement, len(resultM2))
	}

	return res, err
}

func (uc *orderUsecase) ReportHotspotMonth(c context.Context) (res domain.OrderHotspotMonthReport, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res.MonthBreakdown = 0
	res.MonthCurrent = 0
	res.Winning = "-"
	res.WinningData = 0
	res.OtherData = 0

	year := time.Now().Year()
	month := int(time.Now().Month())

	resArr, err := uc.Repository.FindMonthCreated(ctx, month, year)
	if err != nil {
		logrus.Error(err)
	}

	resultM := slices.DeleteFunc(resArr, func(i domain.Order) bool { return i.OrderHotspot.IDBatch == "" })

	res.MonthCurrent = len(resultM)

	if month == 1 {
		month = 12
		year = year - 1

		resArr2, err := uc.Repository.FindMonthCreated(ctx, month, year)
		if err != nil {
			logrus.Error(err)
		}

		resultM2 := slices.DeleteFunc(resArr2, func(i domain.Order) bool { return i.OrderHotspot.IDBatch == "" })

		res.MonthBreakdown = len(resultM2)
	} else {

		month = month - 1
		resArr2, err := uc.Repository.FindMonthCreated(ctx, month, year)
		if err != nil {
			logrus.Error(err)
		}

		resultM2 := slices.DeleteFunc(resArr2, func(i domain.Order) bool { return i.OrderHotspot.IDBatch == "" })

		res.MonthBreakdown = len(resultM2)
	}

	groupByIDPackage := helper.GroupByProperty(resultM, func(o domain.Order) uint {
		return o.OrderHotspot.IDPackage
	})

	for _, group := range groupByIDPackage {
		if len(group) > res.WinningData {
			res.WinningData = len(group)
			res.Winning = group[0].OrderHotspot.Package.Name
			res.OtherData = (len(resultM) - len(group))
		}
	}

	return res, err
}

func (uc *orderUsecase) Find(c context.Context, lastID uint, limit int) (res []domain.Order, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Paginate(ctx, lastID, limit)
	if err != nil {
		logrus.Error(err)
	}
	// order := domain.Order{}
	// res, err = uc.Repository.Find(ctx, &order)
	// if err != nil {
	// 	logrus.Error(err)
	// }

	return
}

func (uc *orderUsecase) Client(ctx context.Context, request *domain.OrderClientRequest) (err error) {
	c, cancel := context.WithTimeout(ctx, uc.Timeout)
	defer cancel()

	// convert client amount
	amountUint64, err := strconv.ParseUint(request.Amount, 10, 64)
	if err != nil {
		logrus.Error(err)
		return err
	}

	clt := domain.Client{ID: request.IDClient}
	client, err := uc.RClient.First(&clt)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if amountUint64 != uint64(client.Session) {
		return domain.ErrMinimumAmountRequired
	}

	// check if client order keep active
	resoc, err := uc.Repository.FindOrderClient(ctx, &domain.OrderClient{IDClient: request.IDClient})
	if err != nil {
		logrus.Error(err)
		return err
	}

	// delete radcheck
	err = uc.RRadcheck.Delete(c, &domain.Rdcheck{Username: client.Username})
	if err != nil {
		logrus.Error(err)
		return err
	}

	// delete radreply
	err = uc.RRadreply.Delete(c, &domain.Radreply{Username: client.Username})
	if err != nil {
		logrus.Error(err)
		return err
	}

	payments := []domain.Payment{}

	layoutFormat := "02 Jan 2006 15:04:05"
	now := time.Now()
	if len(resoc) > 0 {
		wita, err := time.LoadLocation("Asia/Makassar")
		if err != nil {
			logrus.Error(err)
		}

		for _, rs := range resoc {
			fmt.Println("resoc.expire", rs.Expire)
		}
		now, err = time.ParseInLocation(layoutFormat, resoc[0].Expire, wita)
		if err != nil {
			logrus.Error(err)
		}
	}

	if client.ValidityUnit == "DAY" {
		now = now.AddDate(0, 0, int(client.ValidityValue))
	} else if client.ValidityUnit == "MONTH" {
		now = now.AddDate(0, int(client.ValidityValue), 0)
	} else {
		now = now.Add(time.Hour * time.Duration(client.ValidityValue))
	}

	radchecks := []domain.Rdcheck{}
	radreplys := []domain.Radreply{}

	var timeUnix string

	for {
		// cari kode batch yang unik
		timeUnix = strconv.FormatInt(time.Now().Unix(), 10)
		_, err := uc.Repository.First(c, &domain.Order{IDBatch: timeUnix})
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error(err)
			return err
		}

		paymentKredit := domain.Payment{
			IDBatch: timeUnix,
			Status:  "kredit",
			Amount:  float64(client.Price),
		}
		payments = append(payments, paymentKredit)
		if request.Payment == "cash" {
			paymentDebit := domain.Payment{
				IDBatch: timeUnix,
				Status:  "debit",
				Amount:  float64(client.Price),
			}

			payments = append(payments, paymentDebit)
		}
		// bypass ip macaddress
		total := 0
		for _, value := range client.ClientBinding {
			macAddress := strings.ToUpper(value.Mac)
			connectionName := "conn->" + macAddress
			packetName := "pack->" + macAddress
			// bypass ke mikrotik
			// result, err := uc.RIPBinding.Finds(&domain.IPBinding{MacAddress: macAddress})
			for _, ros := range uc.ROSClient {
				result, err := uc.RIPBinding.Finds(ros.Client, &domain.IPBinding{MacAddress: macAddress})
				if err != nil {
					logrus.Error(err)
					return err
				}
				updatedMac := 0
				for ix := 0; ix < len(result); ix++ {
					if result[ix].MacAddress == macAddress {
						// update
						result[ix].Type = "bypassed"
						result[ix].Comment = now.Format(layoutFormat)
						uc.RIPBinding.Updates(ros.Client, &result[ix])

						updatedMac += 1
					}
				}
				if updatedMac == 0 {
					uc.RIPBinding.Saves(ros.Client, &domain.IPBinding{
						MacAddress: macAddress,
						Comment:    now.Format(layoutFormat),
						Type:       "bypassed",
					})
				}

				resultFirewall, err := uc.RFirewall.Firsts(ros.Client, &domain.Firewall{SrcMacAddress: macAddress})
				if err != nil {
					logrus.Error(err)
					return err
				}

				if resultFirewall != (domain.Firewall{}) {
					resultFirewall.Comment = now.Format(layoutFormat)
					resultFirewall.NewConnectionMark = connectionName
					err = uc.RFirewall.Updates(ros.Client, &resultFirewall)
					if err != nil {
						logrus.Error(err)
						return err
					}
				} else {
					uc.RFirewall.Saves(ros.Client, &domain.Firewall{
						Chain:             "prerouting",
						Action:            "mark-connection",
						NewConnectionMark: connectionName,
						Passthrough:       "yes",
						SrcMacAddress:     macAddress,
						Comment:           now.Format(layoutFormat),
					})
				}

				resultFirewallPacket, err := uc.RFirewall.Firsts(ros.Client, &domain.Firewall{ConnectionMark: connectionName})
				if err != nil {
					logrus.Error(err)
					return err
				}

				if resultFirewallPacket != (domain.Firewall{}) {
					resultFirewallPacket.NewPacketMark = packetName
					resultFirewallPacket.Comment = now.Format(layoutFormat)
					err = uc.RFirewall.Updates(ros.Client, &resultFirewallPacket)
					if err != nil {
						logrus.Error(err)
						return err
					}
				} else {
					err = uc.RFirewall.Saves(ros.Client, &domain.Firewall{
						Chain:          "prerouting",
						Action:         "mark-packet",
						ConnectionMark: connectionName,
						Passthrough:    "no",
						NewPacketMark:  packetName,
						Comment:        now.Format(layoutFormat),
					})

					if err != nil {
						logrus.Error(err)
						return err
					}
				}

				resultSimpleQueue, err := uc.RSimpleQueue.Firsts(ros.Client, &domain.SimpleQueue{Name: packetName})
				if err != nil {
					logrus.Error(err)
					return err
				}

				if resultSimpleQueue != (domain.SimpleQueue{}) {
					resultSimpleQueue.MaxLimit = client.Speed
					resultSimpleQueue.PacketMarks = packetName
					resultSimpleQueue.Target = "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
					err = uc.RSimpleQueue.Updates(ros.Client, &resultSimpleQueue)
					if err != nil {
						logrus.Error(err)
						return err
					}
				} else {
					err = uc.RSimpleQueue.Saves(ros.Client, &domain.SimpleQueue{
						Name:        packetName,
						PacketMarks: packetName,
						MaxLimit:    client.Speed,
						Target:      "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16",
					})
					if err != nil {
						logrus.Error(err)
						return err
					}
				}

			}

			// resultFirewall, err := uc.RFirewall.First(&domain.Firewall{SrcMacAddress: macAddress})
			// if err != nil {
			// 	logrus.Error(err)
			// 	return err
			// }

			// if resultFirewall != (domain.Firewall{}) {
			// 	resultFirewall.Comment = now.Format(layoutFormat)
			// 	resultFirewall.NewConnectionMark = connectionName
			// 	err = uc.RFirewall.Update(&resultFirewall)
			// 	if err != nil {
			// 		logrus.Error(err)
			// 		return err
			// 	}
			// } else {
			// 	uc.RFirewall.Save(&domain.Firewall{
			// 		Chain:             "prerouting",
			// 		Action:            "mark-connection",
			// 		NewConnectionMark: connectionName,
			// 		Passthrough:       "yes",
			// 		SrcMacAddress:     macAddress,
			// 		Comment:           now.Format(layoutFormat),
			// 	})
			// }

			// resultFirewallPacket, err := uc.RFirewall.First(&domain.Firewall{ConnectionMark: connectionName})
			// if err != nil {
			// 	logrus.Error(err)
			// 	return err
			// }

			// if resultFirewallPacket != (domain.Firewall{}) {
			// 	resultFirewallPacket.NewPacketMark = packetName
			// 	resultFirewallPacket.Comment = now.Format(layoutFormat)
			// 	err = uc.RFirewall.Update(&resultFirewallPacket)
			// 	if err != nil {
			// 		logrus.Error(err)
			// 		return err
			// 	}
			// } else {
			// 	err = uc.RFirewall.Save(&domain.Firewall{
			// 		Chain:          "prerouting",
			// 		Action:         "mark-packet",
			// 		ConnectionMark: connectionName,
			// 		Passthrough:    "no",
			// 		NewPacketMark:  packetName,
			// 		Comment:        now.Format(layoutFormat),
			// 	})

			// 	if err != nil {
			// 		logrus.Error(err)
			// 		return err
			// 	}
			// }

			// resultSimpleQueue, err := uc.RSimpleQueue.First(&domain.SimpleQueue{Name: packetName})
			// if err != nil {
			// 	logrus.Error(err)
			// 	return err
			// }

			// if resultSimpleQueue != (domain.SimpleQueue{}) {
			// 	resultSimpleQueue.MaxLimit = client.Speed
			// 	resultSimpleQueue.PacketMarks = packetName
			// 	resultSimpleQueue.Target = "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
			// 	err = uc.RSimpleQueue.Update(&resultSimpleQueue)
			// 	if err != nil {
			// 		logrus.Error(err)
			// 		return err
			// 	}
			// } else {
			// 	err = uc.RSimpleQueue.Save(&domain.SimpleQueue{
			// 		Name:        packetName,
			// 		PacketMarks: packetName,
			// 		MaxLimit:    client.Speed,
			// 		Target:      "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16",
			// 	})
			// 	if err != nil {
			// 		logrus.Error(err)
			// 		return err
			// 	}
			// }

			if total == int(client.Session) {
				break
			}
			total++
		}

		// jika session dan total sama maka tak perlu membuat radcheck dan radreply
		if len(client.ClientBinding) < int(client.Session) {
			sisa := int(client.Session) - len(client.ClientBinding)
			if math.Signbit(float64(sisa)) == false && sisa != 0 {
				for i := 0; i < sisa; i++ {
					// buat radcheck
					username := client.Username + strconv.Itoa((i + 1))
					radchecks = append(radchecks, domain.Rdcheck{Username: username, Attribute: "Cleartext-Password", OP: ":=", Value: client.Password})
					radchecks = append(radchecks, domain.Rdcheck{Username: username, Attribute: "Simultaneous-Use", OP: ":=", Value: "1"})
					radchecks = append(radchecks, domain.Rdcheck{Username: username, Attribute: "Expiration", OP: ":=", Value: now.Format(layoutFormat)})
					// buat radreply
					radreplys = append(radreplys, domain.Radreply{Username: username, Attribute: "Mikrotik-Rate-Limit", OP: "=", Value: client.Speed})
				}
			}
		}

		err = uc.Repository.Client(c, &domain.Order{
			IDUser:     request.IDUser,
			IDCustomer: client.IDCustomer,
			IDBatch:    timeUnix,
			Amount:     uint(amountUint64),
		}, &domain.OrderClient{
			IDBatch:   timeUnix,
			IDClient:  client.ID,
			Expire:    now.Format(layoutFormat),
			OrderDate: time.Now(),
		}, radchecks, radreplys, payments)

		if err != nil {
			logrus.Error(err)
		}

		return err
	}
}

func (uc *orderUsecase) Order(ctx context.Context, request *domain.OrderRequest) (err error) {
	c, cancel := context.WithTimeout(ctx, uc.Timeout)
	defer cancel()

	amountInt64, err := strconv.ParseUint(request.Amount, 10, 64)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// cari paket yang dimaksud
	pkg, err := uc.RPackage.First(request.IDPackage)
	if err != nil {
		logrus.Error(err)
		return err
	}

	_, err = uc.RCustomer.First(c, &domain.Customer{ID: request.IDCustomer})
	if err != nil {
		logrus.Error(err)
		return err
	}

	payments := []domain.Payment{}

	layoutFormat := "02 Jan 2006 15:04:05"
	now := time.Now()
	if request.Activated {
		if pkg.ValidityUnit == "DAY" {
			now = now.AddDate(0, 0, int(pkg.ValidityValue))
		} else if pkg.ValidityUnit == "MONTH" {
			now = now.AddDate(0, int(pkg.ValidityValue), 0)
		} else {
			now = now.Add(time.Hour * time.Duration(pkg.ValidityValue))
		}
	}
	vouchers := []domain.Vcr{}
	radchecks := []domain.Rdcheck{}

	var timeUnix string

	min := 4
	max := 8
	config := generator.Config{
		Length:                     uint(rand.Intn(max-min) + min),
		IncludeSymbols:             request.IncludeSymbols,
		IncludeNumbers:             request.IncludeNumbers,
		IncludeLowercaseLetters:    request.IncludeLowercaseLetters,
		IncludeUppercaseLetters:    request.IncludeUppercaseLetters,
		ExcludeSimilarCharacters:   request.ExcludeSimilarCharacters,
		ExcludeAmbiguousCharacters: request.ExcludeAmbiguousCharacters,
	}
	g, _ := generator.New(&config)

	for {
		// cari kode batch yang unik
		timeUnix = strconv.FormatInt(time.Now().Unix(), 10)
		_, err := uc.Repository.First(c, &domain.Order{IDBatch: timeUnix})
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error(err)
			return err
		}

		paymentKredit := domain.Payment{
			IDBatch: timeUnix,
			Status:  "kredit",
			Amount:  float64(pkg.Price) * float64(amountInt64),
		}
		payments = append(payments, paymentKredit)
		if request.Payment == "cash" {
			paymentDebit := domain.Payment{
				IDBatch: timeUnix,
				Status:  "debit",
				Amount:  float64(pkg.Price) * float64(amountInt64),
			}

			payments = append(payments, paymentDebit)
		}

		for {
			// cari kode voucher yang unik
			pwd, err := g.Generate()
			if err != nil {
				logrus.Error(err)
				return err
			}

			// VOUCHER
			_, err = uc.RVcr.First(c, &domain.Vcr{Username: *pwd})
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Error(err)
				return err
			}

			// RADCHECK
			_, err = uc.RRadcheck.First(c, &domain.Rdcheck{Username: *pwd})
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Error(err)
				return err
			}

			if len(vouchers) < int(amountInt64) {
				vouchers = append(vouchers, domain.Vcr{Username: *pwd, Pkg: pkg.ID, Batchcode: timeUnix})
				radchecks = append(radchecks, domain.Rdcheck{Username: *pwd, Attribute: "Cleartext-Password", OP: ":=", Value: *pwd})
				radchecks = append(radchecks, domain.Rdcheck{Username: *pwd, Attribute: "User-Profile", OP: ":=", Value: pkg.Profile})
				if request.Activated {
					// kalkulasi expire
					radchecks = append(radchecks, domain.Rdcheck{Username: *pwd, Attribute: "Expiration", OP: ":=", Value: now.Format(layoutFormat)})
				}
			} else {
				err = uc.Repository.Batch(
					c,
					&domain.Order{
						IDBatch:    timeUnix,
						IDUser:     request.IDUser,
						IDCustomer: request.IDCustomer,
						Amount:     uint(amountInt64),
					},
					&domain.OrderHotspot{
						IDBatch:                    timeUnix,
						IncludeSymbols:             request.IncludeSymbols,
						IncludeNumbers:             request.IncludeNumbers,
						IncludeLowercaseLetters:    request.IncludeLowercaseLetters,
						IncludeUppercaseLetters:    request.IncludeUppercaseLetters,
						ExcludeSimilarCharacters:   request.ExcludeSimilarCharacters,
						ExcludeAmbiguousCharacters: request.ExcludeAmbiguousCharacters,
						Activated:                  request.Activated,
						IDPackage:                  request.IDPackage,
					},
					vouchers,
					radchecks,
					payments,
				)
				if err != nil {
					logrus.Error(err)
					return err
				}
				return nil
			}
		}

	}
}
