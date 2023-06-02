package api

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

var (
	allowedContentType = []string{"image/jpg", "image/jpeg", "image/webp", "image/png"}
)

type SliderService interface {
	ReturnSlider() (string, error)
	SaveSlider(*multipart.Form, *fiber.Ctx) error
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

func (s Slider) createSlider(ctx *fiber.Ctx) error {

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	title := form.Value["title"]
	if title == nil || len(title[0]) < 4 || len(title[0]) > 300 {
		s.log.Debugw("createSlider", "form.Value title", "title is blank or out of range limits")
		ctx.JSON(structs.SetResponse(http.StatusBadRequest, "title is blank or out of range limits"))

		return fiber.NewError(fiber.StatusBadRequest, "title is blank or out of range limits")
	}

	file := form.File["thumb"]
	if file != nil {
		s.log.Debugw("createSlider", "file-name", file[0].Filename, "file-size", file[0].Size)
	}
	if file == nil || !isAllowedContentType(allowedContentType, file[0].Header["Content-Type"][0]) {
		s.log.Debugw("createSlider", "form.File", err.Error())
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, err.Error()))

		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	size := file[0].Size
	if size > 5*1024*1024 {
		s.log.Debugw("createSlider", "form.File", "file too large")
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, "file too large"))

		return fiber.NewError(http.StatusInternalServerError, "file too large")
	}

	if err := s.Service.SaveSlider(form, ctx); err != nil {
		ctx.JSON(structs.SetResponse(fiber.StatusInternalServerError, err.Error()))
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(structs.SetResponse(http.StatusOK, "success"))
}

func isAllowedContentType(allowedList []string, contentType string) bool {
	for _, i := range allowedList {
		if i == contentType {
			return true
		}
	}

	return false
}
