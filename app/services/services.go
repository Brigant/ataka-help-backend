package services

import (
	"strings"

	"github.com/google/uuid"
)

const (
	uploadDirectory = "static/uploads/"
	filePermition   = 0o666
)

type RepoInterface interface {
	CardsRepo
	PartnersRepo
	SliderRepo
	ContactRepo
}

type Services struct {
	CardsService
	PartnersService
	SliderService
	ReportService
	ContactService
	FeedbackService
}

func NewService(repo RepoInterface) Services {
	return Services{
		CardsService{Repo: repo},
		PartnersService{Repo: repo},
		SliderService{Repo: repo},
		ReportService{},
		ContactService{Repo: repo},
		FeedbackService{},
	}
}

func uniqueFilePath(fileName, path string) string {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]

	uniqueFileName := uuid.New().String() + "." + fileExt

	result := path + uniqueFileName

	return result
}
