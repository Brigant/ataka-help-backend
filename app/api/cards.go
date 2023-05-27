package api

import (
	"fmt"
	"net/http"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

const uploadDirectory = "static/uploads/"

type CardService interface {
	ReturnCards() (string, error)
	SaveCard(structs.Card) error
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
	result, err := h.Service.ReturnCards()
	if err != nil {
		return fmt.Errorf("cannot ReturnCarsd: %w", err)
	}

	h.log.Infow("TEST", "val", result)

	if err := ctx.SendString(result); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}

func (h CardHandler) createCard(ctx *fiber.Ctx) error {
	ctx.Accepts("image/png/webp")

	title := ctx.FormValue("title")

	if len(title) < 4 || len(title) > 300 {
		h.log.Debugw("createCard", "form.Value title", "bad title")
		ctx.JSON(structs.SetResponse(http.StatusBadRequest, "bad title"))

		return nil
	}

	description := ctx.FormValue("description")
	if len(description) < 1 {
		h.log.Debugw("createCard", "form.Value description", "bad description")
		ctx.JSON(structs.SetResponse(http.StatusBadRequest, "bad description"))

		return nil
	}

	alt := ctx.FormValue("alt")
	if len(alt) < 1 {
		h.log.Debugw("createCard", "form.Value alt", "bad description")
		ctx.JSON(structs.SetResponse(http.StatusBadRequest, "bad description"))

		return nil
	}
	file, err := ctx.FormFile("thumb")
	h.log.Debugw("createCard", "file-name", file.Filename, "file-size", file.Size)
	if err != nil {
		h.log.Debugw("createCard", "form.File", err.Error())
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, err.Error()))

		return nil
	}

	if file.Size > 5000000 {
		h.log.Debugw("createCard", "form.File", "file to large")
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, "file to large"))

		return nil
	}

	// card := structs.NewCard(title, uploadDirectory+file.Filename, alt, description)
	card := structs.Card{
		Title:       title,
		Thumb:       uploadDirectory + file.Filename,
		Alt:         alt,
		Description: description,
	}

	h.Service.SaveCard(card)

	if err := ctx.SaveFile(file, "static/uploads/"+file.Filename); err != nil {
		h.log.Errorw("createCard", "SaveFile", err.Error())
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, err.Error()))

		return nil
	}

	if err := ctx.JSON(structs.SetResponse(http.StatusOK, "success")); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
