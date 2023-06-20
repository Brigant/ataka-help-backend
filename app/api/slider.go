package api

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

const (
	maxFileSize = 5 * 1024 * 1024
)

type SliderService interface {
	ReturnSlider(context.Context) ([]structs.Slide, error)
	SaveSlider(context.Context, *multipart.Form) error
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
	ctxUser := ctx.UserContext()

	ctxWithDeadline, cancel := context.WithDeadline(ctxUser, time.Now().Add(2*time.Second))

	defer cancel()

	response, err := s.Service.ReturnSlider(ctxWithDeadline)
	if err != nil {
		s.log.Errorw("getSlider", "getSlider error", err.Error())

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	result := structs.SliderResponse{
		Code:   fiber.StatusOK,
		Slider: response,
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (s Slider) createSlider(ctx *fiber.Ctx) error {
	allowedFileExtentions := []string{"jpg", "jpeg", "webp", "png"}

	const (
		minTitle = 4
		maxTitle = 300
		minAlt   = 10
		maxAlt   = 30
	)

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	title := form.Value["title"]

	if title == nil || len(title[0]) < minTitle || len(title[0]) > maxTitle {
		s.log.Debugw("createSlider", "form.Value title", "title is blank or out of range limits")

		return fiber.NewError(fiber.StatusBadRequest, "title is blank or out of range limits")
	}

	file := form.File["thumb"]

	if file != nil {
		s.log.Debugw("createSlider", "file-name", file[0].Filename, "file-size", file[0].Size)
	}

	if file == nil || !isAllowedFileExtention(allowedFileExtentions, file[0].Filename) {
		s.log.Debugw("createSlider", "form.File", err.Error())

		return fiber.NewError(fiber.StatusBadRequest, "thumb is absent")
	}

	alt := form.Value["alt"]

	if alt == nil || len(alt[0]) < minAlt || len(alt[0]) > maxAlt {
		s.log.Debugw("createSlider", "form.Value alt", "alt is blank or out of limits")

		return fiber.NewError(fiber.StatusBadRequest, "alt is blank or out of limits")
	}

	size := file[0].Size

	if size > maxFileSize {
		s.log.Debugw("createSlider", "form.File", "file too large")

		return fiber.NewError(fiber.StatusBadRequest, "file too large")
	}

	if err := s.Service.SaveSlider(ctx.Context(), form); err != nil {
		s.log.Errorw("createSlider", "createSlider error", err.Error())

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(structs.SetResponse(http.StatusOK, "success"))
}
