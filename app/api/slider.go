package api

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/gofiber/fiber/v2"
)

type SliderService interface {
	ReturnSlider() (string, error)
}

type Slider struct {
	Service SliderService
	log     *logger.Logger
}

func NewSliderHandler(service SliderService, log *logger.Logger) Slider {
	return Slider{
		Service: service,
		log:     log,
	}
}

func (s Slider) getSlider(ctx *fiber.Ctx) error {
	result, err := s.Service.ReturnSlider()
	if err != nil {
		return fmt.Errorf("cannot ReturnSlider: %w", err)
	}

	s.log.Infow("TEST", "val", result)

	if err = ctx.SendString(result); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
