package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type CardsRepo interface {
	SelectAllCards(context.Context, structs.CardQueryParameters) ([]structs.Card, error)
	InsertCard(context.Context, structs.Card) error
	CountRowsTable(context.Context, string) (int, error)
	SelectCardByID(context.Context, string) (structs.Card, error)
	DelCardByID(context.Context, string) error
}

type CardsService struct {
	Repo CardsRepo
}

func (s CardsService) ReturnCards(ctx context.Context, params structs.CardQueryParameters) ([]structs.Card, int, error) {
	cards, err := s.Repo.SelectAllCards(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("error happens while SelectAllCards: %w", err)
	}

	total, err := s.Repo.CountRowsTable(ctx, "cards")
	if err != nil {
		return nil, 0, fmt.Errorf("error happens while CountRowsTable: %w", err)
	}

	return cards, total, nil
}

func (s CardsService) SaveCard(ctx context.Context, form *multipart.Form) error {
	fileHeader := form.File["thumb"][0]

	card := structs.Card{
		Title:       form.Value["title"][0],
		Thumb:       uniqueFilePath(fileHeader.Filename, uploadDirectory),
		Alt:         form.Value["alt"][0],
		Description: []byte(form.Value["description"][0]),
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("error happens while fileHeader.Open(): %w", err)
	}

	osFile, err := os.OpenFile(card.Thumb, os.O_WRONLY|os.O_CREATE, filePermition)
	if err != nil {
		return fmt.Errorf("error happens while os.OpenFile(): %w", err)
	}

	defer osFile.Close()

	written, err := io.Copy(osFile, file)
	if err != nil {
		return fmt.Errorf(" written bytes: %v, error happens while io.Copy(): %w", written, err)
	}

	if err := s.Repo.InsertCard(ctx, card); err != nil {
		if err := os.Remove(card.Thumb); err != nil {
			return fmt.Errorf("error happens while remove file: %w", err)
		}

		return fmt.Errorf("error happens while InsertCard: %w", err)
	}

	return nil
}

func (s CardsService) ReturnCardByID(ctx context.Context, id string) (structs.Card, error) {
	card, err := s.Repo.SelectCardByID(ctx, id)
	if err != nil {
		return structs.Card{}, fmt.Errorf("error while select card: %w", err)
	}

	return card, nil
}

func (s CardsService) DeleteCardByID(ctx context.Context, id string) error {
	if err := s.Repo.DelCardByID(ctx, id); err != nil {
		return fmt.Errorf("error while delete card: %w", err)
	}

	return nil
}
