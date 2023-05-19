package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Card
	Partner
	logger.Logger
}

func NewHandler(services ServiceInterfaces, log logger.Logger) Handler {
	return Handler{
		Card{Service: services},
		Partner{Service: services},
		log,
	}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Get("/", h.Card.getCards)
	app.Get("/partners", h.Partner.Get)
}
