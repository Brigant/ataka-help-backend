package api

import "github.com/gofiber/fiber/v2"

type CardsHandler struct {
	CardService CardService
}

func NewCardsHandler(service CardService) CardsHandler {
	return CardsHandler{
		CardService: service,
	}
}

func (h CardsHandler) getCards(c *fiber.Ctx) error {
	result, err := h.CardService.ReturnCards(c)
	if err != nil {
		return err
	}

	return c.SendString(result)
}
