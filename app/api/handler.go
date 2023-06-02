package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/logger"
)

var (
	allowedContentType = []string{"image/jpg", "image/jpeg", "image/webp", "image/png"}
)

type ServiceInterfaces interface {
	CardService
	PartnerService
	SliderService
}

type Handler struct {
	Card    CardHandler
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

func isAllowedContentType(allowedList []string, contentType string) bool {
	for _, i := range allowedList {
		if i == contentType {
			return true
		}
	}

	return false
}
