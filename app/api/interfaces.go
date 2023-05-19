package api

type ServiceInterfaces struct {
	CardService
	PartnerService
}

type CardService interface {
	ReturnCards() (string, error)
}

type PartnerService interface {
	GetAll() (string, error)
}
