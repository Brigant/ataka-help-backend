package services

import (
	"fmt"
	"strings"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/google/uuid"
)

const (
	uploadDirectory = "static/uploads/"
	filePermition   = 0o666
	templatPath     = "./app/services/template/emailTemplate.gohtml"
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

func NewService(repo RepoInterface, cfg config.SMTP) (Services, error) {
	feedbackService, err := NewFeedbackService(cfg)
	if err != nil {
		return Services{}, fmt.Errorf("error in NewFeedbackService(): %w", err)
	}

	return Services{
		CardsService{Repo: repo},
		PartnersService{Repo: repo},
		SliderService{Repo: repo},
		ReportService{},
		ContactService{Repo: repo},
		feedbackService,
	}, nil
}

func uniqueFilePath(fileName, path string) string {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]

	uniqueFileName := uuid.New().String() + "." + fileExt

	result := path + uniqueFileName

	return result
}
