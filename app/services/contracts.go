package services

import "github.com/gofiber/fiber/v2"

type CardsStorage interface {
	SelectAllCards(c *fiber.Ctx) (string, error)
}
