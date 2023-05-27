package services

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type CardsRepo interface {
	SelectAllCards() (string, error)
	InsertCard(structs.Card) error
}

type CardsService struct {
	Repo CardsRepo
}

func (s CardsService) ReturnCards() (string, error) {
	result, err := s.Repo.SelectAllCards()
	if err != nil {
		return "", fmt.Errorf("error happens while select: %w", err)
	}

	return result, nil
}

func (s CardsService) SaveCard(card structs.Card) error {
	if err := s.Repo.InsertCard(card); err != nil {
		return fmt.Errorf("error happens while select: %w", err)
	}

	return nil
}
