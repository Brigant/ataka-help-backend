package services

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

const uploadDirectory = "static/uploads/"

type CardsRepo interface {
	SelectAllCards(structs.CardQueryParameters, context.Context) ([]structs.Card, error)
	InsertCard(structs.Card, context.Context) error
	CountRowsTable(string, context.Context) (int, error)
}

type CardsService struct {
	Repo CardsRepo
}

func (s CardsService) ReturnCards(params structs.CardQueryParameters, ctx context.Context) ([]structs.Card, int, error) {
	cards, err := s.Repo.SelectAllCards(params, ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error happens while SelectAllCards: %w", err)
	}

	total, err := s.Repo.CountRowsTable("cards", ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error happens while CountRowsTable: %w", err)
	}

	return cards, total, nil
}

func (s CardsService) SaveCard(form *multipart.Form, ctx *fiber.Ctx) error {
	file := form.File["thumb"][0]

	descr, err := json.Marshal(form.Value["description"][0])
	if err != nil {
		return fmt.Errorf("error happens while Marshal description: %w", err)
	}

	card := structs.Card{
		Title:       form.Value["title"][0],
		Thumb:       uniqueFilePath(file.Filename, uploadDirectory),
		Alt:         form.Value["alt"][0],
		Description: descr,
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
