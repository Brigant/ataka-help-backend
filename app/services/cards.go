package services

import (
	"fmt"
)

type CardsService struct {
	Repo CardsRepository
}

func (s CardsService) ReturnCards() (string, error) {
	result, err := s.Repo.SelectAllCards()
	if err != nil {
		return "", fmt.Errorf("error happens while select: %w", err)
	}

	return result, nil
}
