package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/logger"
)

const (
	fileLimit     = 2 * 1024 * 1024
	defaultLimit  = 6
	defaultOffset = 0
)

type ServiceInterfaces interface {
	CardService
	PartnerService
	ReportService
	ContactService
	FeedbackService
}

type Handler struct {
	Card     CardHandler
	Partner  Partner
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
		Feedback: NewFeedbackHandler(services, log),
	}
}

func isAllowedContentType(allowedList []string, contentType string) bool {
	for _, i := range allowedList {
		if i == contentType {
			return true
		}
	}

	return false
}
