package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type usersUsecase struct {
	Timeout    time.Duration
	Repository domain.UsersRepository
}

// NewUsecase ...
func NewUsecase(t time.Duration, r domain.UsersRepository) domain.UsersUsecase {
	return &usersUsecase{Timeout: t, Repository: r}
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(dbPass string, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass)) == nil
}

func (uc *usersUsecase) Register(c context.Context, nik string, telegramID string) (res string, err error) {
	return "", domain.ErrBadParamInput
}

func (uc *usersUsecase) ImportUsers(c context.Context, nik []string) (err error) {
	return domain.ErrBadParamInput
}

func (uc *usersUsecase) LevelFirst(level *domain.Level) (res domain.Level, err error) {
	res, err = uc.Repository.LevelFirst(level)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *usersUsecase) LevelFind(level *domain.Level) (res []domain.Level, err error) {
	res, err = uc.Repository.LevelFind(level)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *usersUsecase) LevelSave(request *domain.LevelRequest) (err error) {
	err = uc.Repository.LevelSave(&domain.Level{Name: request.Name})
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *usersUsecase) LevelDelete(level *domain.Level) (err error) {
	err = uc.Repository.LevelDelete(level)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *usersUsecase) Find(user *domain.Usr) (res []domain.Usr, err error) {
	res, err = uc.Repository.Find(user)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *usersUsecase) Save(user *domain.Usr) (err error) {
	level, err := uc.Repository.LevelFirst(&domain.Level{ID: user.IDLevel})
	if err != nil {
		logrus.Error(err)
		return
	}

	user.Password = hashAndSalt([]byte(user.Password))
	user.UserLevel = level
	err = uc.Repository.Save(user)
	if err != nil {
		logrus.Error(err)
		return
	}

	// ok, err := initializers.ENFORCER.AddGroupingPolicy(fmt.Sprint(user.ID), level.Name)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }

	// if !ok {
	// 	err = fmt.Errorf("pengguna gagal di daftarkan pada grup")
	// 	return
	// }

	return
}

func (uc *usersUsecase) ResetPassword(request *domain.UsrResetRequest) (err error) {
	// check if user find and
	user, err := uc.Repository.First(&domain.Usr{ID: request.ID})
	if err != nil {
		logrus.Error(err)

		return
	}

	dPassword := "12345678"
	user.Password = hashAndSalt([]byte(dPassword))
	user.CreatedAt = time.Now()

	err = uc.Repository.Save(&user)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *usersUsecase) Delete(request *domain.Usr) (err error) {
	err = uc.Repository.Delete(request)
	if err != nil {
		logrus.Error(err)
	}

	return
}

// Login ...
func (uc *usersUsecase) Login(c context.Context, username string, password string) (res domain.JWTCustomClaims, err error) {
	_, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	usr := domain.Usr{}
	usr.Email = username

	resUser, err := uc.Repository.First(&usr)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = fmt.Errorf("email atau kata sandi salah")
		return
	} else if (err != nil) && (!errors.Is(err, gorm.ErrRecordNotFound)) {
		logrus.Error(err)
		return
	}

	if resUser == (domain.Usr{}) {
		err = fmt.Errorf("email atau kata sandi salah")
		logrus.Error(err)
		return
	}

	if !comparePasswords(resUser.Password, password) {
		err = fmt.Errorf("email atau kata sandi salah")
		logrus.Error(err)
		return
	}

	claims := domain.JWTCustomClaims{
		Username: resUser.Email,
		ID:       resUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(viper.GetString(`server.secret`)))
	if err != nil {
		logrus.Error(err)
		return
	}

	claims.Token = t

	return claims, nil
}
