package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Card struct {
	Service CardService
}

func NewCardsHandler(service CardService) Card {
	return Card{
		Service: service,
	}
}

func (h Card) getCards(ctx *fiber.Ctx) error {
	result, err := h.Service.ReturnCards()
	if err != nil {
		return fmt.Errorf("cannot ReturnCarsd: %w", err)
	}

	if err := ctx.SendString(result); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
