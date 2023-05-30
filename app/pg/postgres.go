package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	ErrCodeUniqueViolation     = "unique_violation"
	ErrCodeNoData              = "no_data"
	ErrCodeForeignKeyViolation = "foreign_key_violation"
	ErrCodeUndefinedColumn     = "undefined_column"
)

type Repo struct {
	db sqlx.DB
}

// Returns an object of the Ropository.
func NewRepository(cfg config.Config) (Repo, error) {
	database, err := sqlx.Connect("postgres", fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Database, cfg.DB.Password, cfg.DB.SSLmode))
	if err != nil {
		return Repo{}, fmt.Errorf("cannot connect to db: %w", err)
	}

	return Repo{db: *database}, nil
}

func (r Repo) Close() error {
	return r.db.Close()
}

func (r Repo) SelectAllCards(offest, limit int, ctx context.Context) ([]structs.Card, error) {
	query := `
		SELECT * FROM public.cards c
		ORDER BY c.created DESC
		Limit $1
		OFFSET $2;
	`
	cards := []structs.Card{}

	if err := r.db.SelectContext(ctx, cards, query, limit, offest); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, structs.ErrNotFound
		}

		return nil, fmt.Errorf("an error occurs while getting the cards list: %w", err)
	}

	return cards, nil
}

func (r Repo) InsertCard(card structs.Card, ctx context.Context) error {
	const expectedEffectedRow = 1

	query := `INSERT INTO public.cards
	(title, thumb, alt, description)
	VALUES(:title, :thumb, :alt, :description);`

	result, err := r.db.NamedExecContext(ctx, query, card)
	if err != nil {
		pqError := new(pq.Error)
		if errors.As(err, &pqError) && pqError.Code.Name() == ErrCodeUniqueViolation {
			return structs.ErrUniqueRestriction
		}

		return fmt.Errorf("error in NamedEx: %w", err)
	}

	effectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("the error is in RowsAffected: %w", err)
	}

	if effectedRows != expectedEffectedRow {
		return structs.ErrDatabaseInserting
	}

	return nil
}

func (r Repo) SelectAllPartners() (string, error) {
	return "some partners from db", nil
}
