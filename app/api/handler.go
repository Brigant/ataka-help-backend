package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/logger"
)

const (
	fileLimit     = 5 * 1024 * 1024
	defaultLimit  = 6
	defaultOffset = 0
)

type ServiceInterfaces interface {
	CardService
	PartnerService
	ReportService
}

type Handler struct {
	Card    CardHandler
	Partner Partner
	Report  ReportHandler
}

func NewHandler(services ServiceInterfaces, log *logger.Logger) Handler {
	return Handler{
		Card:    NewCardsHandler(services, log),
		Partner: NewParnerHandler(services, log),
		Report:  NewReportHandler(services, log),
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
