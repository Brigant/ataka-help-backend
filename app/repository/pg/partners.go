package pg

import "github.com/jmoiron/sqlx"

type PartnersRepo struct {
	db *sqlx.DB
}

func (p PartnersRepo) SelectAllPartners() (string, error) {
	return "some partners from db", nil
}
