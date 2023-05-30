package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

const uploadDirectory = "static/uploads/"

type CardsRepo interface {
	SelectAllCards(offest, limit int, ctx context.Context) ([]structs.Card, error)
	InsertCard(structs.Card, context.Context) error
}

type CardsService struct {
	Repo CardsRepo
}

func (s CardsService) ReturnCards(offset, limit int, ctx context.Context) ([]structs.Card, int, error) {
	cards, err := s.Repo.SelectAllCards(6, 1, ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error happens while SelectAllCards: %w", err)
	}

	return cards, 0, nil
}

func (s CardsService) SaveCard(form *multipart.Form, ctx *fiber.Ctx) error {
	file := form.File["thumb"][0]

	card := structs.Card{
		Title:       form.Value["title"][0],
		Thumb:       uniqueFilePath(file.Filename, uploadDirectory),
		Alt:         form.Value["alt"][0],
		Description: form.Value["description"][0],
	}

	if err := ctx.SaveFile(file, card.Thumb); err != nil {
		return fmt.Errorf("error happens while SaveFile: %w", err)
	}

	if err := s.Repo.InsertCard(card, ctx.Context()); err != nil {
		if err := os.Remove(card.Thumb); err != nil {
			return fmt.Errorf("error happens while remove file: %w", err)
		}

		return fmt.Errorf("error happens while InsertCard: %w", err)
	}

	return nil
}
