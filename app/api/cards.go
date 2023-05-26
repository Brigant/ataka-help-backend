package api

import (
	"fmt"
	"net/http"

	"github.com/baza-trainee/ataka-help-backend/app/core"
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/gofiber/fiber/v2"
)

type CardService interface {
	ReturnCards() (string, error)
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

	if form, err := ctx.MultipartForm(); err == nil {
		// => *multipart.Form

		if description := form.Value["description"]; len(description) > 0 {
			// Get key value:
			fmt.Println(description)
		}

		// Get all files from "documents" key:
		f := form.File["image"]
		// => []*multipart.FileHeader
		for _, file := range f {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"
		}
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.log.Errorf("TEST", "val", err.Error())
		return fmt.Errorf("cannot read image: %w", err)
	}

	card := core.Card{
		Title: ctx.FormValue("Title"),
		Alt:   ctx.FormValue("Alt"),
		Image: "static/uploads/" + file.Filename,
	}

	ctx.SaveFile(file, card.Image)
	fmt.Printf("%+v\n", card)
	if err := ctx.JSON(core.SetResponse(http.StatusOK, "success")); err != nil {
		return fmt.Errorf("some err: %w", err)
	}

	return nil
}
