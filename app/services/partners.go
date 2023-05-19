package services

type PartnersService struct {
	Repo PartnersRepository
}

func (s PartnersService) GetAll() (string, error) {
	str, err := s.Repo.SelectAllPartners()
	if err != nil {
		return "", err
	}
	return str, err
}
