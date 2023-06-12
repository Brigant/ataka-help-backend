package api

import (
	"strings"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
)

var allowedContentType = []string{"image/jpg", "image/jpeg", "image/webp", "image/png"}

const (
	fileLimit     = 2 * 1024 * 1024
	defaultLimit  = 6
	defaultOffset = 0
)

type ServiceInterfaces interface {
	CardService
	PartnerService
	SliderService
	ReportService
	ContactService
	FeedbackService
}

type Handler struct {
	Card     CardHandler
	Partner  Partner
	Slider   Slider
	Report   ReportHandler
	Contact  ContactHandler
	Feedback FeedbackHandler
}

func NewHandler(services ServiceInterfaces, log *logger.Logger) Handler {
	return Handler{
		Card:     NewCardsHandler(services, log),
		Partner:  NewParnerHandler(services, log),
		Report:   NewReportHandler(services, log),
		Contact:  NewContactHandler(services, log),
		Slider:   NewSliderHandler(services, log),
		Feedback: NewFeedbackHandler(services, log),
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
