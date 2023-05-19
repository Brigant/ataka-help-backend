package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Partner struct {
	Service PartnerService
}

func (h Partner) Get(ctx *fiber.Ctx) error {
	str, err := h.Service.GetAll()
	if err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	if err := ctx.SendString(str); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
