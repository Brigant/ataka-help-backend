package services

import (
	"fmt"
)

type SliderRepo interface {
	SelectSlider() (string, error)
}

type SliderService struct {
	Repo SliderRepo
}

func (s SliderService) ReturnSlider() (string, error) {
	result, err := s.Repo.SelectSlider()
	if err != nil {
		return "", fmt.Errorf("error happens while select: %w", err)
	}

	return result, nil
}
