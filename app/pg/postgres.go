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
	return r.db.Close()
}

func (r Repo) SelectAllCards() (string, error) {
	return "some string from DB", nil
}

func (r Repo) SelectAllPartners() (string, error) {
	return "some partners from db", nil
}

func (r Repo) SelectSlider() (string, error) {
	return "array of slider images from db", nil
}

func (r Repo) InsertSlider(slider structs.Slider, ctx context.Context) error {

	query := `INSERT INTO public.slider (title, thumb, alt, description)
			  VALUES(:title, :thumb);`

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
