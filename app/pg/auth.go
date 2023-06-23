package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/lib/pq"
)

func (r Repo) FindEmailWithPasword(ctx context.Context, identity structs.IdentityData) (string, error) {
	var userID string

	query := `SELECT id FROM public.users WHERE email=$1 and password=$2 `

	select {
	case <-ctx.Done():
		return "", structs.ErrTimeout
	default:
		if err := r.db.GetContext(ctx, &userID, query, identity.Login, identity.Password); err != nil {
			pqErr := new(pq.Error)
			if errors.As(err, &pqErr) && pqErr.Code.Name() == "query_canceled" {
				return "", structs.ErrTimeout
			}

			if errors.Is(err, sql.ErrNoRows) {
				return "", structs.ErrNotFound
			}

			return "", fmt.Errorf("error while GetContext(): %w", err)
		}

		return userID, nil
	}
}

func (r Repo) CheckUserIDWithPasword(ctx context.Context, userID, password string) error {
	query := `SELECT id FROM public.users WHERE id=$1 and password=$2 `

	select {
	case <-ctx.Done():
		return structs.ErrTimeout
	default:
		id := ""
		if err := r.db.GetContext(ctx, &id, query, userID, password); err != nil {
			pqErr := new(pq.Error)
			if errors.As(err, &pqErr) && pqErr.Code.Name() == "query_canceled" {
				return structs.ErrTimeout
			}

			if id == "" {
				return structs.ErrNotFound
			}

			return fmt.Errorf("error while GetContext(): %w", err)
		}

		return nil
	}
}

func (r Repo) UpdatePassord(ctx context.Context, userID, password string) error {
	query := `UPDATE  public.users SET password=$2 WHERE id=$1 `

	select {
	case <-ctx.Done():
		return structs.ErrTimeout
	default:
		if result, err := r.db.ExecContext(ctx, query, userID, password); err != nil {
			pqErr := new(pq.Error)
			if errors.As(err, &pqErr) && pqErr.Code.Name() == "query_canceled" {
				return structs.ErrTimeout
			}

			affected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("error in RowsAffected(): %w", err)
			}

			if affected != expectedAffectedRow {
				return structs.ErrNoRowAffected
			}

			return fmt.Errorf("error while GetContext(): %w", err)
		}

		return nil
	}
}
