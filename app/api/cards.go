package api

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardService interface {
	ReturnCards(context.Context, structs.CardQueryParameters) ([]structs.Card, int, error)
	SaveCard(context.Context, *multipart.Form) error
	ReturnCardByID(context.Context, string) (structs.Card, error)
	DeleteCardByID(context.Context, string) error
}

type CardHandler struct {
	Service CardService
	log     *logger.Logger
}

func NewCardsHandler(service CardService, log *logger.Logger) CardHandler {
	return CardHandler{
		Service: service,
		log:     log,
	}
}

func (h CardHandler) getCards(ctx *fiber.Ctx) error {
	params := structs.CardQueryParameters{}

	if err := ctx.QueryParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	cards, total, err := h.Service.ReturnCards(ctx.UserContext(), params)
	if err != nil && !errors.Is(err, structs.ErrNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := structs.CardsResponse{
		Code:  fiber.StatusOK,
		Total: total,
		Cards: cards,
	}

	return ctx.Status(fiber.StatusOK).JSON(response) // nolint
}

// nolint: cyclop
func (h CardHandler) createCard(ctx *fiber.Ctx) error {
	allowedFileExtentions := []string{"jpg", "jpeg", "webp", "png"}

	const (
		limitNumberItemsFile = 1
		minAltItems          = 10
		maxAltItems          = 30
		minTitle             = 4
		maxTitle             = 300
		minDescription       = 3
	)

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if len(form.File["thumb"]) < limitNumberItemsFile {
		h.log.Debugw("createCard", "form.File", "no thumb was attached")

		return fiber.NewError(fiber.StatusBadRequest, "no thumb was attached")
	}

	fileHeader := form.File["thumb"][0]

	switch {
	case fileHeader == nil || fileHeader.Size > fileLimit || !isAllowedFileExtention(allowedFileExtentions, fileHeader.Filename):
		h.log.Debugw("createCard", "form.File", "required thumb not biger then 5 Mb and format jpg/jpeg/webp")

		return fiber.NewError(fiber.StatusBadRequest, "required thumb not bigger then 5 Mb and format jpg/jpeg/webp")
	case form.Value["title"] == nil || len(form.Value["title"][0]) < minTitle || len(form.Value["title"][0]) > maxTitle:
		h.log.Debugw("createCard", "form.Vlaues", "required title more than 3 letters and less than 300")

		return fiber.NewError(fiber.StatusBadRequest, "required title more than 3 letters and less than 300")
	case form.Value["alt"] == nil || len(form.Value["alt"][0]) < minAltItems || len(form.Value["alt"][0]) > maxAltItems:
		h.log.Debugw("createCard", "form.Vlaues", "required alt")

		return fiber.NewError(fiber.StatusBadRequest, "required alt")
	case form.Value["description"] == nil || len(form.Value["description"][0]) < minDescription:
		h.log.Debugw("createCard", "form.Vlaues", "required description")

		return fiber.NewError(fiber.StatusBadRequest, "required description")
	}

	if err := h.Service.SaveCard(ctx.Context(), form); err != nil {
		if errors.Is(err, structs.ErrUniqueRestriction) {
			h.log.Errorw("SaveCard", "error", err.Error())

			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(structs.SetResponse(fiber.StatusCreated, "success")) // nolint
}

func (h CardHandler) findCard(ctx *fiber.Ctx) error {
	param := struct {
		ID string `params:"id"`
	}{}

	if err := ctx.ParamsParser(&param); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_, err := uuid.Parse(param.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id is not uuid type")
	}

	card, err := h.Service.ReturnCardByID(ctx.Context(), param.ID)
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(card)
}

func (h CardHandler) deleteCard(ctx *fiber.Ctx) error {
	param := struct {
		ID string `params:"id"`
	}{}

	if err := ctx.ParamsParser(&param); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_, err := uuid.Parse(param.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "id is not uuid type")
	}

	if err := h.Service.DeleteCardByID(ctx.Context(), param.ID); err != nil {
		if errors.Is(err, structs.ErrNoRowAffected) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success"))
}
