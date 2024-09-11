package usecase

import (
	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type clientUsecase struct {
	Repository    domain.ClientRepository
	PkgRepository domain.PackageRepository
}

func NewUsecase(repo domain.ClientRepository, pkg domain.PackageRepository) domain.ClientUsecase {
	return &clientUsecase{Repository: repo, PkgRepository: pkg}
}

func (uc *clientUsecase) Find(client *domain.Client) (res []domain.Client, err error) {

	res, err = uc.Repository.Find(client)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *clientUsecase) FindBinding(binding *domain.ClientBinding) (res []domain.ClientBinding, err error) {

	res, err = uc.Repository.FindBinding(binding)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *clientUsecase) First(id uint) (res domain.Client, err error) {
	client := domain.Client{
		ID: id,
	}

	res, err = uc.Repository.First(&client)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *clientUsecase) Save(client *domain.Client) (err error) {
	err = uc.Repository.Save(client)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *clientUsecase) SaveBinding(binding *domain.ClientBinding) (err error) {
	err = uc.Repository.SaveBinding(binding)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func (uc *clientUsecase) Delete(client *domain.ClientRequest) (err error) {
	// first get valid id
	cli := domain.Client{ID: client.ID}
	_, err = uc.Repository.First(&cli)
	if err != nil {
		logrus.Error(err)
		return
	}

	err = uc.Repository.Delete(&cli)
	if err != nil {
		logrus.Error()
	}

	return
}

func (uc *clientUsecase) DeleteBinding(binding *domain.ClientBinding) (err error) {
	// first get valid id
	bdi := domain.ClientBinding{ID: binding.ID}

	err = uc.Repository.DeleteBinding(&bdi)
	if err != nil {
		logrus.Error()
	}

	return
}
