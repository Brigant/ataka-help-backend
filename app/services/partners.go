package services

import "fmt"

type PartnersRepo interface {
	SelectAllPartners() (string, error)
}

type PartnersService struct {
	Repo PartnersRepo
}

func (s PartnersService) ReturnPartners() (string, error) {
	str, err := s.Repo.SelectAllPartners()
	if err != nil {
		return "", fmt.Errorf("error in PartnerService GetAll: %w", err)
	}

	return str, nil
}
