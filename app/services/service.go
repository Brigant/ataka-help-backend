package services

import "github.com/baza-trainee/ataka-help-backend/app/repository/pg"

type Services struct {
	CardsService
	PartnersService
}

func NewService(repo pg.Repository) Services {
	return Services{
		CardsService{Repo: repo.CardsRepo},
		PartnersService{Repo: repo.PartnersRepo},
	}
}
