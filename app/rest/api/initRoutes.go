package api

import (
	"github.com/baza-trainee/ataka-help-backend/config"
	"github.com/gofiber/fiber/v2"
)

// The structure describes the dependencies.
type Deps struct {
	CardsService CardService
}

type Handler struct {
	Cards CardsHandler
}

func NewHandler(deps Deps) Handler {
	return Handler{
		Cards: NewCardsHandler(deps.CardsService),
	}
}

func (h *Handler) InitRoutes(cfg config.Config) *fiber.App {
	// todo: make fiber config

	app := fiber.New()

	app.Get("/", h.Cards.getCards)

	return app
}
