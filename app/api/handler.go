package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Card    Card
	Partner Partner
}

func NewHandler(services ServiceInterfaces, log *logger.Logger) Handler {
	return Handler{
		Card:    NewCardsHandler(services, log),
		Partner: NewParnerHandler(services, log),
	}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Get("/", h.Card.getCards)
	app.Get("/partners", h.Partner.Get)
}
