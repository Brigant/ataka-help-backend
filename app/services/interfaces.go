package services

type RepoInterface interface {
	CardsRepository
	PartnersRepository
}

type CardsRepository interface {
	SelectAllCards() (string, error)
}

type PartnersRepository interface {
	SelectAllPartners() (string, error)
}
