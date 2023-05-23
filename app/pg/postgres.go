package pg

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // nececarry blank import
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

func (r Repo) SelectAllCards() (string, error) {
	return "some string from DB", nil
}

func (r Repo) SelectAllPartners() (string, error) {
	return "some partners from db", nil
}

func (r Repo) Close() error {
	return r.db.Close()
}
