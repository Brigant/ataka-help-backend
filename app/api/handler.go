package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/logger"
)

type ServiceInterfaces interface {
	CardService
	PartnerService
}

type Handler struct {
	Card    CardHandler
	Partner Partner
}

func NewHandler(services ServiceInterfaces, log *logger.Logger) Handler {
	return Handler{
		Card:    NewCardsHandler(services, log),
		Partner: NewParnerHandler(services, log),
	}
}

