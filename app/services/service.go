package services

type Services struct {
	CardsService
	PartnersService
}

func NewService(repo RepoInterface) Services {
	return Services{
		CardsService{Repo: repo},
		PartnersService{Repo: repo},
	}
}
