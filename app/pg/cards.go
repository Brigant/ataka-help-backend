package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/lib/pq"
)

func (r Repo) SelectAllCards(ctx context.Context, params structs.CardQueryParameters) ([]structs.Card, error) {
	cards := []structs.Card{}
	var limit *int = &params.Limit

	query := `
		SELECT id, title, thumb, alt, description, created, modified
		FROM public.cards 
		ORDER BY created DESC
		Limit $1 OFFSET $2;
	`
	if params.Limit == 0 {
		limit = nil
	}

	select {
	case <-ctx.Done():
		return nil, structs.ErrTimeout
	default:
		if err := r.db.SelectContext(ctx, &cards, query, limit, params.Offset); err != nil {
			pqErr := new(pq.Error)
			if errors.As(err, &pqErr) && pqErr.Code.Name() == "query_canceled" {
				return nil, structs.ErrTimeout
			}

			return nil, fmt.Errorf("an error occurs while QueryContext: %w", err)
		}

		return cards, nil
	}
}

func (r Repo) InsertCard(ctx context.Context, card structs.Card) error {
	query := `INSERT INTO public.cards
	(title, thumb, alt, description)
	VALUES($1, $2, $3, $4::json);`

	select {
	case <-ctx.Done():
		return structs.ErrTimeout
	default:
		result, err := r.db.ExecContext(ctx, query, card.Title, card.Thumb, card.Alt, card.Description)
		if err != nil {
			pqError := new(pq.Error)
			if errors.As(err, &pqError) && pqError.Code.Name() == "query_canceled" {
				return structs.ErrTimeout
			}

			if errors.As(err, &pqError) && pqError.Code.Name() == ErrCodeUniqueViolation {
				return structs.ErrUniqueRestriction
			}

			return fmt.Errorf("error in NamedEx: %w", err)
		}

		effectedRows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("the error is in RowsAffected: %w", err)
		}

		if effectedRows != expectedAffectedRow {
			return structs.ErrNoRowAffected
		}

		return nil
	}
}

func (r Repo) SelectCardByID(ctx context.Context, id string) (structs.Card, error) {
	query := `SELECT * FROM public.cards WHERE id=$1`

	card := structs.Card{}

	if err := r.db.GetContext(ctx, &card, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.Card{}, structs.ErrNotFound
		}

		return structs.Card{}, fmt.Errorf("error in GetContext(): %w", err)
	}

	select {
	case <-ctx.Done():
		return card, structs.ErrTimeout
	default:
		if err := r.db.GetContext(ctx, &card, query, id); err != nil {
			pqErr := new(pq.Error)
			if errors.As(err, &pqErr) && pqErr.Code.Name() == "query_canceled" {
				return structs.Card{}, structs.ErrTimeout
			}

			if errors.Is(err, sql.ErrNoRows) {
				return structs.Card{}, structs.ErrNotFound
			}

			return structs.Card{}, fmt.Errorf("error in GetContext(): %w", err)
		}

		return card, nil
	}
}

func (r Repo) DelCardByID(ctx context.Context, id string) error {
	query := `DELETE FROM public.cards WHERE id=$1`

	sqlResult, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error in ExecContext(): %w", err)
	}

	affectedRows, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("the error is in RowsAffected: %w", err)
	}

	if affectedRows != expectedAffectedRow {
		return structs.ErrNoRowAffected
	}

	select {
	case <-ctx.Done():
		return structs.ErrTimeout
	default:
		sqlResult, err := r.db.ExecContext(ctx, query, id)
		if err != nil {
			pqErr := new(pq.Error)
			if errors.As(err, &pqErr) && pqErr.Code.Name() == "query_canceled" {
				return structs.ErrTimeout
			}

			return fmt.Errorf("error in ExecContext(): %w", err)
		}

		affectedRows, err := sqlResult.RowsAffected()
		if err != nil {
			return fmt.Errorf("the error is in RowsAffected: %w", err)
		}

		if affectedRows != expectedAffectedRow {
			return structs.ErrNoRowAffected
		}

		return nil
	}
}
