package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

const (
	uploadDirectory = "static/uploads/"
)

type SliderService interface {
	ReturnSlider() (string, error)
	SaveSlider(slider structs.Slider) error
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
	// validating file format
	image := ctx.Accepts("image/png/jpg/webp")
	if len(image) < 1 {
		s.log.Errorw("createSlider", "form.File", "invalid image format")
		ctx.JSON(structs.SetResponse(http.StatusUnsupportedMediaType, "invalid image format"))
	}

	// validating title size
	title := ctx.FormValue("title")
	if len(title) < 4 || len(title) > 300 {
		s.log.Debugw("createSlider", "form.Value title", "bad title")
		ctx.JSON(structs.SetResponse(http.StatusBadRequest, "bad title"))

		return errors.New("invalid title size")
	}

	// getting file
	file, err := ctx.FormFile("thumb")
	s.log.Debugw("createSlider", "file-name", file.Filename, "file-size", file.Size)
	if err != nil {
		s.log.Debugw("createSlider", "form.File", err.Error())
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, err.Error()))

		return errors.New(err.Error())
	}

	// validating file size
	if file.Size > 5000000 {
		s.log.Debugw("createSlider", "form.File", "file too large")
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, "file too large"))

		return errors.New("file too large")
	}

	// creating the object to transport it to --> service --> repository
	slider := structs.Slider{
		Title: title,
		Thumb: uploadDirectory + file.Filename,
	}

	// transport it to --> service
	s.Service.SaveSlider(slider)

	// save file
	if err := ctx.SaveFile(file, uploadDirectory+file.Filename); err != nil {
		s.log.Errorw("createSlider", "SaveFile", err.Error())
		ctx.JSON(structs.SetResponse(http.StatusInternalServerError, err.Error()))

		return errors.New(err.Error())
	}

	// return success code
	if err := ctx.JSON(structs.SetResponse(http.StatusOK, "success")); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
