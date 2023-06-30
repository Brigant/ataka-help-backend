package api

import (
	"strings"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/logger"
)

const (
	fileLimit     = 2 * 1024 * 1024
	defaultLimit  = 0
	defaultOffset = 0
	apiVersion1 = "/api/v1"
)

type ServiceInterfaces interface {
	CardService
	PartnerService
	SliderService
	ReportService
	ContactService
	FeedbackService
	AutService
}

type Handler struct {
	Card     CardHandler
	Partner  Partner
	Slider   Slider
	Report   ReportHandler
	Contact  ContactHandler
	Feedback FeedbackHandler
	Auth     AuthHandler
}

func NewHandler(services ServiceInterfaces, log *logger.Logger, cfg config.Config) Handler {
	return Handler{
		Card:     NewCardsHandler(services, log),
		Partner:  NewPartnerHandler(services, log, cfg.Server),
		Report:   NewReportHandler(services, log),
		Contact:  NewContactHandler(services, log),
		Slider:   NewSliderHandler(services, log, cfg.Server),
		Feedback: NewFeedbackHandler(services, log),
		Auth:     NewAuthHandler(services, log, cfg.Auth),
	}
}

func isAllowedFileExtention(allowedList []string, fileName string) bool {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]
	for _, i := range allowedList {
		if i == fileExt {
			return true
		}
	}

	return false
}

func symbolsCounter(sentence string) int {
	runes := []rune(sentence)

	return len(runes)
}
