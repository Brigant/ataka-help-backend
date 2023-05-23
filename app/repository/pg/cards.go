package pg

import (
	"github.com/jmoiron/sqlx"
)

type CardsRepo struct {
	db *sqlx.DB
}

func (db CardsRepo) SelectAllCards() (string, error) {
	return "some string from DB", nil
}
