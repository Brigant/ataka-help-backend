package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CardsHandler struct {
	CardService CardService
}

func NewCardsHandler(service CardService) CardsHandler {
	return CardsHandler{
		CardService: service,
	}
}

func (h CardsHandler) getCards(ctx *fiber.Ctx) error {
	result, err := h.CardService.ReturnCards(ctx)
	if err != nil {
		return fmt.Errorf("cannot ReturnCarsd: %w", err)
	}

	if err := ctx.SendString(result); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
