package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type PartnersRepo interface {
	SelectAllPartners(context.Context, structs.PartnerQueryParameters) ([]structs.Partner, error)
	InsertPartner(context.Context, structs.Partner, chan struct{}) error
	DelPartnerByID(context.Context, string) (string, error)
	CountRowsTable(context.Context, string) (int, error)
}

type PartnersService struct {
	Repo PartnersRepo
}

func (s PartnersService) ReturnPartners(ctx context.Context, params structs.PartnerQueryParameters) ([]structs.Partner, int, error) {
	partners, err := s.Repo.SelectAllPartners(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("error happens while SelectAllPartners: %w", err)
	}

	total, err := s.Repo.CountRowsTable(ctx, "partners")
	if err != nil {
		return nil, 0, fmt.Errorf("error happens while CountRowsTable: %w", err)
	}

	return partners, total, nil
}

func (p PartnersService) SavePartner(ctx context.Context, form *multipart.Form, chWell chan struct{}) error {
	file := form.File["logo"][0]

	partner := structs.Partner{
		Alt:  form.Value["alt"][0],
		Logo: uniqueFilePath(file.Filename, uploadDirectory),
	}

	fileOpened, err := file.Open()
	if err != nil {
		return fmt.Errorf("error happens while file.Open(): %w", err)
	}

	osFile, err := os.OpenFile(partner.Logo, os.O_WRONLY|os.O_CREATE, filePermition)
	if err != nil {
		return fmt.Errorf("error happens while os.OpenFile(): %w", err)
	}

	defer osFile.Close()

	written, err := io.Copy(osFile, fileOpened)
	if err != nil {
		return fmt.Errorf(" written bytes: %v, error happens while io.Copy(): %w", written, err)
	}

	if err := p.Repo.InsertPartner(ctx, partner, chWell); err != nil {
		if err := os.Remove(partner.Logo); err != nil {
			return fmt.Errorf("error happens while remove file: %w", err)
		}

		return fmt.Errorf("error happens while inserting: %w", err)
	}

	return nil
}

func (p PartnersService) DeletePartnerByID(ctx context.Context, ID string, chWell chan struct{}) error {
	objectPath, err := p.Repo.DelPartnerByID(ctx, ID)
	if err != nil {
		return fmt.Errorf("error while delete partner: %w", err)
	}

	if err := os.Remove(objectPath); err != nil {
		return fmt.Errorf("error happens while remove file: %w", err)
	}

	chWell <- struct{}{}

	close(chWell)

	return nil
}
