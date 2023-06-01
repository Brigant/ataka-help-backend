package api

import (
	"context"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

type ContactService interface {
	Modify(context.Context, structs.Contact) error
	Obtain(context.Context) (structs.Contact, error)
}

type ContactHandler struct {
	Service ContactService
	log     *logger.Logger
}

func NewContactHandler(service ContactService, log *logger.Logger) ContactHandler {
	return ContactHandler{
		Service: service,
		log:     log,
	}
}

func (h ContactHandler) Edit(ctx *fiber.Ctx) error {
	contact := structs.Contact{}

	if err := ctx.BodyParser(contact); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.Service.Modify(ctx.Context(), contact); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success")) // nolint
}

// TODO: Implement Handler for returning contact
func (c ContactHandler) Get(ctx *fiber.Ctx) error {
	return nil
}
