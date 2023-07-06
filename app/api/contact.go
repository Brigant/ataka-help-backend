package api

import (
	"context"
	"errors"

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

func (h ContactHandler) edit(ctx *fiber.Ctx) error {
	contact := structs.Contact{}

	if err := ctx.BodyParser(&contact); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if contact.Email == "" || contact.Phone1 == "" || contact.Phone2 == "" {
		return fiber.NewError(fiber.StatusBadRequest, "all contact fields are rewuired")
	}

	if err := h.Service.Modify(ctx.Context(), contact); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success")) // nolint
}

func (h ContactHandler) get(ctx *fiber.Ctx) error {
	contact, err := h.Service.Obtain(ctx.Context())
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(contact)
}
