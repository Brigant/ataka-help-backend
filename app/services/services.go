package services

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/google/uuid"
)

const (
	uploadDirectory = "static/uploads/"
	filePermition   = 0o666
	templatPath     = "./app/services/template/emailTemplate.go.html"
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
	template, err := template.ParseFiles(templatPath)
	if err != nil {
		return Services{}, fmt.Errorf("error in checkGoogleCaptcha(): %w", err)
	}

	return Services{
		CardsService{Repo: repo},
		PartnersService{Repo: repo},
		SliderService{Repo: repo},
		ReportService{},
		ContactService{Repo: repo},
		FeedbackService{cfg: cfg, templateFile: template},
	}, nil
}

func uniqueFilePath(fileName, path string) string {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]

	uniqueFileName := uuid.New().String() + "." + fileExt

	result := path + uniqueFileName

	return result
}
