package api

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

type SliderService interface {
	ReturnSlider() ([]structs.Slide, error)
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
	response, err := s.Service.ReturnSlider()
	if err != nil {
		s.log.Errorw("getSlider", "getSLider error", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	result := structs.SliderResponse{
		Code:   fiber.StatusOK,
		Slider: response,
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (s Slider) createSlider(ctx *fiber.Ctx) error {

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	title := form.Value["title"]
	if title == nil || len(title[0]) < 4 || len(title[0]) > 300 {
		s.log.Debugw("createSlider", "form.Value title", "title is blank or out of range limits")

		return fiber.NewError(fiber.StatusBadRequest, "title is blank or out of range limits")
	}

	file := form.File["thumb"]
	if file != nil {
		s.log.Debugw("createSlider", "file-name", file[0].Filename, "file-size", file[0].Size)
	}
	if file == nil || !isAllowedContentType(allowedContentType, file[0].Header["Content-Type"][0]) {
		s.log.Debugw("createSlider", "form.File", err.Error())

		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	size := file[0].Size
	if size > 5*1024*1024 {
		s.log.Debugw("createSlider", "form.File", "file too large")

		return fiber.NewError(http.StatusInternalServerError, "file too large")
	}

	if err := s.Service.SaveSlider(ctx.Context(), form); err != nil {
		s.log.Errorw("createSlider", "createSlider error", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(structs.SetResponse(http.StatusOK, "success"))
}
