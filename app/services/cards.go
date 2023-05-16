package services

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CardsService struct {
	CardsStorage CardsStorage
}

func NewCardsService(storage CardsStorage) CardsService {
	return CardsService{
		CardsStorage: storage,
	}
}

func (s CardsService) ReturnCards(c *fiber.Ctx) (string, error) {
	result, err := s.CardsStorage.SelectAllCards(c)
	if err != nil {
		return "", fmt.Errorf("error happens while select: %w", err)
	}

	return result, nil
}
