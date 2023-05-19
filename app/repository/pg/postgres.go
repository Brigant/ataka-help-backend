package pg

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // nececarry blank import
)

type Repository struct {
	CardsRepo
	PartnersRepo
}

// NewPostgresDB function returns object of datatabase.
func NewPostgresDB(cfg config.Config) (*sqlx.DB, error) {
	database, err := sqlx.Connect("postgres", fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Database, cfg.DB.Password, cfg.DB.SSLmode))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	return database, nil
}

// Returns an object of the Ropository.
func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		CardsRepo:    CardsRepo{db: db},
		PartnersRepo: PartnersRepo{db: db},
	}
}
