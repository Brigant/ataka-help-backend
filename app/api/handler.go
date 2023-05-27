package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/logger"
)

type ServiceInterfaces interface {
	CardService
	PartnerService
	SliderService
}

type Handler struct {
	Card    Card
	Partner Partner
	Slider  Slider
}

func NewHandler(services ServiceInterfaces, log *logger.Logger) Handler {
	return Handler{
		Card:    NewCardsHandler(services, log),
		Partner: NewParnerHandler(services, log),
		Slider:  NewSliderHandler(services, log),
	}
}
