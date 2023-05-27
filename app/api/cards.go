package api

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/gofiber/fiber/v2"
)

type CardService interface {
	ReturnCards() (string, error)
}

type Card struct {
	Service CardService
	log     *logger.Logger
}

func NewCardsHandler(service CardService, log *logger.Logger) Card {
	return Card{
		Service: service,
		log:     log,
	}
}

func (h Card) getCards(ctx *fiber.Ctx) error {
	result, err := h.Service.ReturnCards()
	if err != nil {
		return fmt.Errorf("cannot ReturnCards: %w", err)
	}

	h.log.Infow("TEST", "val", result)

	if err := ctx.SendString(result); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
