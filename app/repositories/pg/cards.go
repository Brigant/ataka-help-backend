package pg

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type CardsDB struct {
	db *sqlx.DB
}

func NewAccountDB(db *sqlx.DB) CardsDB {
	return CardsDB{
		db: db,
	}
}

func (s CardsDB) SelectAllCards(ctx *fiber.Ctx) (string, error) {
	return "some string from DB", nil
}
