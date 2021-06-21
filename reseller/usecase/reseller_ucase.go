package usecase

import (
	"context"
	"math"
	"os/exec"
	"strings"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
)

type resellerUsecase struct {
	Timeout               time.Duration
	Repository            domain.ResellerRepository
	PackageRepository     domain.PackageRepository
	RadcheckRepository    domain.RadcheckRepository
	TransactionRepository domain.TransactionRepository
	RadacctRepository     domain.RadacctRepository
}

func NewUsecase(timeout time.Duration, r domain.ResellerRepository, pr domain.PackageRepository, rcr domain.RadcheckRepository, tr domain.TransactionRepository, rra domain.RadacctRepository) domain.ResellerUsecase {
	return &resellerUsecase{Timeout: timeout, Repository: r, PackageRepository: pr, RadcheckRepository: rcr, TransactionRepository: tr, RadacctRepository: rra}
}

func (u *resellerUsecase) ChangeProfile(c context.Context, voucher string, profile string) (res string, err error) {
	ctx, cancel := context.WithTimeout(c, u.Timeout)
	defer cancel()

	attr := "User-Profile"

	radcheck := domain.Radcheck{}
	radcheck.Username = &voucher
	radcheck.Attribute = &attr

	rcArr, err := u.RadcheckRepository.Get(ctx, radcheck)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if len(rcArr) <= 0 {
		logrus.Error("item profile not found: ", voucher)
		return "", domain.ErrNotFound
	}

	if *rcArr[0].Value == profile {
		logrus.Error("voucher in profile")
		return "", domain.ErrConflict
	}

	// check apakah memungkinkan untuk di update dengan melihat apakah dia mempunyai masa aktiv
	// jika telah expired maka tolak semua pembaharuan
	// jika aktiv maka lakukan update profile

	wita, err := time.LoadLocation("Asia/Makassar")

	if err != nil {
		logrus.Error(err)
		return "", domain.ErrInternalServerError
	}

	layoutFormat := "02 Jan 2006 15:04:05"

	attr = "Expiration"

	rcExpiration := domain.Radcheck{}
	rcExpiration.Username = &voucher
	rcExpiration.Attribute = &attr

	rcExpArr, err := u.RadcheckRepository.Get(ctx, radcheck)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if len(rcExpArr) <= 0 {
		logrus.Error("item expiration not found: ", voucher)
		return "", domain.ErrNotFound
	}

	date, err := time.ParseInLocation(layoutFormat, *rcExpArr[0].Value, wita)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	duration := time.Now().Sub(date)
	// if not have duration, kita tidak dapat melanjutkan
	if math.Signbit(duration.Seconds()) == false {
		logrus.Error("expired voucher cannot process")
		return "", domain.ErrInternalServerError
	}

	// ok everithing ok lets change package
	rcUpdate := domain.Radcheck{}
	rcUpdate.ID = rcArr[0].ID
	rcUpdate.Username = rcArr[0].Username
	rcUpdate.Attribute = rcArr[0].Attribute
	rcUpdate.OP = rcArr[0].OP
	rcUpdate.Value = &profile

	err = u.RadcheckRepository.Update(ctx, rcUpdate)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	lRadacct := domain.Radacct{}
	lRadacct.Username = rcArr[0].Username
	radacct, err := u.RadacctRepository.Get(ctx, lRadacct)
	if err != nil {
		logrus.Error(err)
	} else if len(radacct) <= 0 {
		logrus.Error("user not found in table radaccr: ", *rcArr[0].Username)
	} else {
		cmd := exec.Command("sh", "-c", `echo "Acct-Session-Id=`+*radacct[0].Acctsessionid+`,User-Name=`+*radacct[0].Username+`,NAS-IP-Address=`+*radacct[0].Nasipaddress+`,Framed-IP-Address=`+*radacct[0].Framedipaddress+`" | radclient -x `+*radacct[0].Nasipaddress+`:3799 disconnect '`+*radacct[0].Secret+`'`)
		errNotSend := cmd.Run()
		if errNotSend != nil {
			logrus.Error(err)
		}
	}

	res = *rcExpArr[0].Value

	return
}

func (u *resellerUsecase) ChangePackage(c context.Context, voucher string, profile string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.Timeout)
	defer cancel()

	attr := "User-Profile"

	radcheck := domain.Radcheck{}
	radcheck.Username = &voucher
	radcheck.Attribute = &attr

	rcArr, err := u.RadcheckRepository.Get(ctx, radcheck)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if len(rcArr) <= 0 {
		logrus.Error("item profile not found: ", voucher)
		return domain.ErrNotFound
	}

	if *rcArr[0].Value == profile {
		logrus.Error("voucher in profile")
		return domain.ErrConflict
	}

	// check apakah memungkinkan untuk di update dengan melihat apakah dia mempunyai masa aktiv
	// jika telah expired maka tolak semua pembaharuan
	// jika aktiv maka lakukan update profile

	wita, err := time.LoadLocation("Asia/Makassar")

	if err != nil {
		logrus.Error(err)
		return domain.ErrInternalServerError
	}

	layoutFormat := "02 Jan 2006 15:04:05"

	attr = "Expiration"

	rcExpiration := domain.Radcheck{}
	rcExpiration.Username = &voucher
	rcExpiration.Attribute = &attr

	rcExpArr, err := u.RadcheckRepository.Get(ctx, radcheck)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if len(rcExpArr) <= 0 {
		logrus.Error("item expiration not found: ", voucher)
		return domain.ErrNotFound
	}

	date, err := time.ParseInLocation(layoutFormat, *rcExpArr[0].Value, wita)
	if err != nil {
		logrus.Error(err)
		return err
	}

	duration := time.Now().Sub(date)
	// if not have duration, kita tidak dapat melanjutkan
	if math.Signbit(duration.Seconds()) == false {
		logrus.Error("expired voucher cannot process")
		return domain.ErrInternalServerError
	}

	// ok everithing ok lets change package
	rcUpdate := domain.Radcheck{}
	rcUpdate.ID = rcArr[0].ID
	rcUpdate.Username = rcArr[0].Username
	rcUpdate.Attribute = rcArr[0].Attribute
	rcUpdate.OP = rcArr[0].OP
	rcUpdate.Value = &profile

	err = u.RadcheckRepository.Update(ctx, rcUpdate)
	if err != nil {
		logrus.Error(err)
		return err
	}

	lRadacct := domain.Radacct{}
	lRadacct.Username = rcArr[0].Username
	radacct, err := u.RadacctRepository.Get(ctx, lRadacct)
	if err != nil {
		logrus.Error(err)
	} else if len(radacct) <= 0 {
		logrus.Error("user not found in table radaccr: ", *rcArr[0].Username)
	} else {
		cmd := exec.Command("sh", "-c", `echo "Acct-Session-Id=`+*radacct[0].Acctsessionid+`,User-Name=`+*radacct[0].Username+`,NAS-IP-Address=`+*radacct[0].Nasipaddress+`,Framed-IP-Address=`+*radacct[0].Framedipaddress+`" | radclient -x `+*radacct[0].Nasipaddress+`:3799 disconnect '`+*radacct[0].Secret+`'`)
		errNotSend := cmd.Run()
		if errNotSend != nil {
			logrus.Error(err)
		}
	}

	return
}

// Transaction(c context.Context, idUsers int64, idPackage int64) (res Transaction, err error)
func (u *resellerUsecase) Transaction(c context.Context, idUsers int64, idPackage int64) (res domain.Transaction, err error) {
	ctx, cancel := context.WithTimeout(c, u.Timeout)
	defer cancel()

	pac := domain.Package{}
	pac.ID = &idPackage
	packages, err := u.PackageRepository.Get(ctx, pac)
	if err != nil {
		logrus.Error(err)
		return domain.Transaction{}, err
	}

	if len(packages) <= 0 {
		logrus.Error("package not found ", idUsers)
		return domain.Transaction{}, domain.ErrNotFound
	}

	balance, err := u.Balance(ctx, idUsers)
	if err != nil {
		logrus.Error(err)
		return domain.Transaction{}, err
	}

	if balance.Balance < *packages[0].Price {
		logrus.Error("balance not enough :", idUsers)
		return domain.Transaction{}, domain.ErrInSufficientBalance
	}

	safeUsername := ""
	for true {
		// search random password
		username, err := password.Generate(8, 3, 0, true, true)
		if err != nil {
			logrus.Error(err)
			return domain.Transaction{}, err
		}
		radcheck := domain.Radcheck{}
		radcheck.Username = &username

		radchecks, _ := u.RadcheckRepository.Get(ctx, radcheck)
		if len(radchecks) == 0 {
			safeUsername += username
			break
		}
	}

	idTransaction, err := u.Repository.Transaction(ctx, idUsers, safeUsername, safeUsername, packages[0])
	if err != nil {
		logrus.Error(err)
		return domain.Transaction{}, err
	}

	transaction := domain.Transaction{}
	transaction.ID = &idTransaction
	transactions, err := u.TransactionRepository.Get(ctx, transaction)
	if err != nil {
		logrus.Error(err)
		return domain.Transaction{}, err
	}

	if len(transactions) == 0 {
		logrus.Error("transaction not found in: ", idTransaction)
		return domain.Transaction{}, domain.ErrNotFound
	}

	res = transactions[0]

	return
}

// Balance(c context.Context, idUsers int64) (res ResellerBalance, err error)
func (u *resellerUsecase) Balance(c context.Context, idUsers int64) (res domain.ResellerBalance, err error) {
	ctx, cancel := context.WithTimeout(c, u.Timeout)
	defer cancel()

	transaction := domain.Transaction{}
	transaction.IDReseller = &idUsers

	resArr, err := u.TransactionRepository.Get(ctx, transaction)
	if err != nil {
		logrus.Error(err)
		res.Balance = 0
		return res, err
	}

	var in int64 = 0
	var out int64 = 0
	for _, value := range resArr {
		if strings.EqualFold("IN", *value.Status) == true {
			in += *value.Value
			continue
		}
		if strings.EqualFold("OUT", *value.Status) == true {
			out += *value.Value
			continue
		}
	}

	res.Balance = in - out

	return
}
