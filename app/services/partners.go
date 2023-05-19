package services

import "fmt"

type PartnersService struct {
	Repo PartnersRepository
}

func (s PartnersService) GetAll() (string, error) {
	str, err := s.Repo.SelectAllPartners()
	if err != nil {
		return "", fmt.Errorf("error in PartnerService GetAll: %w", err)
	}

	return str, nil
}
