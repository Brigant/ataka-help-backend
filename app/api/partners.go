package api

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/gofiber/fiber/v2"
)

type Partner struct {
	Service PartnerService
	log     *logger.Logger
}

func NewParnerHandler(service PartnerService, log *logger.Logger) Partner {
	return Partner{
		Service: service,
		log:     log,
	}
}

func (h Partner) Get(ctx *fiber.Ctx) error {
	str, err := h.Service.GetAll()
	if err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	h.log.Infow("TEST", "val", str)

	if err := ctx.SendString(str); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
