package api

import "github.com/gofiber/fiber/v2"

type CardService interface {
	ReturnCards(c *fiber.Ctx) (string, error)
}
