package services

type Deps struct {
	CardsStorage CardsStorage
}

type Services struct {
	CardsService CardsService
}

func New(deps Deps) Services {
	return Services{
		CardsService: NewCardsService(deps.CardsStorage),
	}
}
