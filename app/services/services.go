package services

type RepoInterface interface {
	CardsRepo
	PartnersRepo
	SliderRepo
}

type Services struct {
	CardsService
	PartnersService
	SliderService
}

func NewService(repo RepoInterface) Services {
	return Services{
		CardsService{Repo: repo},
		PartnersService{Repo: repo},
		SliderService{Repo: repo},
	}
}
