package services

import (
	"context"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type ContactRepo interface {
	UpdateContact(context.Context, structs.Contact) error
	SelectContact(context.Context) (structs.Contact, error)
}

type ContactService struct {
	Repo ContactRepo
}

func (s ContactService) Modify(ctx context.Context, contact structs.Contact) error {
	if err := s.Repo.UpdateContact(ctx, contact); err!=nil {
		return fmt.Errorf("error occures while UpdateContact(): %w", err)
	}

	return nil
}

func (s ContactService) Obtain(ctxt context.Context) (structs.Contact, error) {
	return structs.Contact{}, nil
}
