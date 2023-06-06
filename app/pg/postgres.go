package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	expectedAffectedRow        = 1
	ErrCodeUniqueViolation     = "unique_violation"
	ErrCodeNoData              = "no_data"
	ErrCodeForeignKeyViolation = "foreign_key_violation"
	ErrCodeUndefinedColumn     = "undefined_column"
)

type Repo struct {
	db *sqlx.DB
}

// Returns an object of the Ropository.
func NewRepository(cfg config.Config) (Repo, error) {
	database, err := sqlx.Connect("postgres", fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Database, cfg.DB.Password, cfg.DB.SSLmode))
	if err != nil {
		return Repo{}, fmt.Errorf("cannot connect to db: %w", err)
	}

	return Repo{db: database}, nil
}

func (r Repo) Close() error {
	return fmt.Errorf("error hapens while db.close: %w", r.db.Close())
}

func (r Repo) SelectAllCards(ctx context.Context, params structs.CardQueryParameters) ([]structs.Card, error) {
	query := `
		SELECT id, title, thumb, alt, description, created, modified
		FROM public.cards c
		ORDER BY c.created DESC
		Limit $1
		OFFSET $2;
	`
	cards := []structs.Card{}

	rows, err := r.db.QueryContext(ctx, query, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("an error occurs while QueryContext: %w", err)
	}

	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurs while rows.Err(): %w", err)
	}

	for rows.Next() {
		card := structs.Card{}

		if err := rows.Scan(
			&card.ID, &card.Title, &card.Thumb, &card.Alt,
			&card.Description, &card.Created, &card.Modified); err != nil {
			return nil, fmt.Errorf("an error occurs while rows.Scan: %w", err)
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (r Repo) InsertCard(ctx context.Context, card structs.Card) error {

	query := `INSERT INTO public.cards
	(title, thumb, alt, description)
	VALUES($1, $2, $3, $4::json);`

	result, err := r.db.ExecContext(ctx, query, card.Title, card.Thumb, card.Alt, card.Description)
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

	if effectedRows != expectedAffectedRow {
		return structs.ErrDatabaseInserting
	}

	return nil
}

func (r Repo) CountRowsTable(ctx context.Context, table string) (int, error) {
	query := `SELECT count(*) as result FROM public.` + table

	var total int

	if err := r.db.GetContext(ctx, &total, query); err != nil {
		return 0, fmt.Errorf("error in GetContext: %w", err)
	}

	return total, nil
}

func (r Repo) SelectAllPartners() (string, error) {
	return "some partners from db", nil
}

func (r Repo) SelectSlider() ([]structs.Slider, error) {
	response := []structs.Slider{}

	query := `SELECT * FROM public.slider AS sld
			  ORDER BY sld.created DESC;`

	err := r.db.Select(&response, query)
	if err != nil {
		return []structs.Slider{}, fmt.Errorf("error happens while slider returning: %w", err)
	}
	return response, nil
}

func (r Repo) InsertSlider(ctx context.Context, slider structs.Slider) error {

	query := `INSERT INTO public.slider (title, thumb)
			  VALUES($1, $2);`

	result, err := r.db.ExecContext(ctx, query, slider.Title, slider.Thumb)
	if err != nil {
		pqError := new(pq.Error)
		if errors.As(err, &pqError) && pqError.Code.Name() == ErrCodeForeignKeyViolation {
			return structs.ErrForeignViolation
		}

		if errors.As(err, &pqError) && pqError.Code.Name() == ErrCodeUniqueViolation {
			return structs.ErrUniqueRestriction
		}

		return fmt.Errorf("error in ExecContext: %w", err)
	}

	effectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("the error is in RowsAffected: %w", err)
	}

	if effectedRows != expectedAffectedRow {
		return structs.ErrDatabaseInserting
	}

	return nil
}
