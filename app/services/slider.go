package services

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type SliderRepo interface {
	SelectSlider() (string, error)
	InsertSlider(slider structs.Slider) error
}

type SliderService struct {
	Repo SliderRepo
}

func (s SliderService) ReturnSlider() (string, error) {
	result, err := s.Repo.SelectSlider()
	if err != nil {
		return "", fmt.Errorf("error happens while selecting: %w", err)
	}

	return result, nil
}

func (s SliderService) SaveSlider(slider structs.Slider) error {
	if err := s.Repo.InsertSlider(slider); err != nil {
		return fmt.Errorf("error happens while inserting: %w", err)
	}

	return nil
}
