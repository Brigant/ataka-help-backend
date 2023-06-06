package services

import (
	"strings"

	"github.com/google/uuid"
)

type RepoInterface interface {
	CardsRepo
	PartnersRepo
	ContactRepo
}

type Services struct {
	CardsService
	PartnersService
	ContactService
}

func NewService(repo RepoInterface) Services {
	return Services{
		CardsService{Repo: repo},
		PartnersService{Repo: repo},
		ContactService{Repo: repo},
	}
}

func uniqueFilePath(fileName, path string) string {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]

	uniqueFileName := uuid.New().String() + "." + fileExt

	result := path + uniqueFileName

	return result
}
