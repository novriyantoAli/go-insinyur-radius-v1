package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		logrus.Error(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
	// 090816
	// $2a$04$p5iYeqT0BZa5iZfllD5IQOhMzZqB4vqitMurDwpWdaOyia0QEJ8DO
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		logrus.Error(err)
		return false
	}

	return true
}

func (uc *usersUsecase) ResetPassword(c context.Context, nik string) (err error) {
	return domain.ErrBadParamInput
}

func (uc *usersUsecase) Register(c context.Context, nik string, telegramID string) (res string, err error) {
	return "", domain.ErrBadParamInput
}

func (uc *usersUsecase) ImportUsers(c context.Context, nik []string) (err error) {
	return domain.ErrBadParamInput
}

// Login ...
func (uc *usersUsecase) Login(c context.Context, username string, password string) (res domain.JWTCustomClaims, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	user := domain.Users{}
	user.Username = &username

	fmt.Println("pass: ", hashAndSalt([]byte(password)))

	resUser, err := uc.Repository.Find(ctx, user)
	if err != nil {
		logrus.Error(err)
		return domain.JWTCustomClaims{}, err
	}

	if resUser == (domain.Users{}) {
		return domain.JWTCustomClaims{}, domain.ErrUsernameOrPasswordWrong
	}

	if !comparePasswords(*resUser.Password, []byte(password)) {
		return domain.JWTCustomClaims{}, domain.ErrUsernameOrPasswordWrong
	}

	claims := domain.JWTCustomClaims{
		Username: *resUser.Username,
		ID:       *resUser.ID,
		Level:    *resUser.Level,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(viper.GetString(`server.secret`)))
	if err != nil {
		logrus.Error(err)
		return domain.JWTCustomClaims{}, err
	}

	claims.Token = &t

	return claims, nil
}
