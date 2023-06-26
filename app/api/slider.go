package api

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	maxFileSize = 5 * 1024 * 1024
	id          = "id"
	nothing     = "nothing"
)

type SliderService interface {
	ReturnSlider(context.Context) ([]structs.Slide, error)
	SaveSlider(context.Context, *multipart.Form, chan struct{}) error
	DeleteSlideByID(context.Context, string, chan struct{}) error
}

type Slider struct {
	Service SliderService
	log     *logger.Logger
	cfg     config.Server
}

func NewSliderHandler(service SliderService, log *logger.Logger, cfg config.Server) Slider {
	return Slider{
		Service: service,
		log:     log,
		cfg:     cfg,
	}
}

func (s Slider) getSlider(ctx *fiber.Ctx) error {
	chErr := make(chan error)

	chWell := make(chan structs.SliderResponse)

	ctxUser := ctx.UserContext()

	ctxWithDeadline, cancel := context.WithDeadline(ctxUser, time.Now().Add(s.cfg.AppWriteTimeout))

	defer cancel()

	go func(chErr chan error, chWell chan structs.SliderResponse) {
		response, err := s.Service.ReturnSlider(ctxWithDeadline)
		if err != nil {
			s.log.Errorw("getSlider", "getSlider error", err.Error())

			chErr <- err

			close(chErr)
		}

		result := structs.SliderResponse{
			Code:   fiber.StatusOK,
			Slider: response,
		}

		chWell <- result

		close(chWell)
	}(chErr, chWell)

	select {
	case result := <-chWell:
		return ctx.Status(fiber.StatusOK).JSON(result)
	case err := <-chErr:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	case <-ctxWithDeadline.Done():
		return ctx.Status(fiber.StatusRequestTimeout).JSON(structs.SetResponse(fiber.StatusRequestTimeout, fiber.ErrRequestTimeout.Message))
	}
}

func (s Slider) createSlide(ctx *fiber.Ctx) error {
	allowedFileExtentions := []string{"jpg", "jpeg", "webp", "png"}

	chErr := make(chan error)

	chWell := make(chan struct{})

	const (
		limitNumberItemsFile = 1
		minTitle             = 4
		maxTitle             = 300
		minAlt               = 10
		maxAlt               = 30
	)

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	title := form.Value["title"]

	if title == nil || symbolsCounter(title[0]) < minTitle || symbolsCounter(title[0]) > maxTitle {
		s.log.Debugw("createSlider", "form.Value title", "title is blank or out of range limits")

		return fiber.NewError(fiber.StatusBadRequest, "title is blank or out of range limits")
	}

	if len(form.File["thumb"]) < limitNumberItemsFile {
		s.log.Debugw("createSlide", "form.File", "thumb is absent")

		return fiber.NewError(fiber.StatusBadRequest, "thumb is absent")
	}

	file := form.File["thumb"]

	if !isAllowedFileExtention(allowedFileExtentions, file[0].Filename) {
		s.log.Debugw("createSlide", "file-name", file[0].Filename, "file-size", form.File["thumb"][0].Size)
	}

	alt := form.Value["alt"]

	if alt == nil || symbolsCounter(alt[0]) < minAlt || symbolsCounter(alt[0]) > maxAlt {
		s.log.Debugw("createSlider", "form.Value alt", "alt is blank or out of limits")

		return fiber.NewError(fiber.StatusBadRequest, "alt is blank or out of limits")
	}

	size := file[0].Size

	if size > maxFileSize {
		s.log.Debugw("createSlide", "form.File", "file too large")

		return fiber.NewError(fiber.StatusBadRequest, "file too large")
	}

	ctxUser := ctx.UserContext()

	ctxWithDeadline, cancel := context.WithDeadline(ctxUser, time.Now().Add(s.cfg.AppWriteTimeout))

	defer cancel()

	go func(chErr chan error, chWell chan struct{}) {
		if err := s.Service.SaveSlider(ctxWithDeadline, form, chWell); err != nil {
			s.log.Errorw("createSlide", "createSlide error", err.Error())

			chErr <- err

			close(chErr)
		}
	}(chErr, chWell)

	select {
	case <-chWell:
		return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success"))
	case err := <-chErr:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	case <-ctxWithDeadline.Done():
		return ctx.Status(fiber.StatusRequestTimeout).JSON(structs.SetResponse(fiber.StatusRequestTimeout, fiber.ErrRequestTimeout.Message))
	}
}

func (s Slider) deleteSlide(ctx *fiber.Ctx) error {
	chErr := make(chan error)

	chWell := make(chan struct{})

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

	ctxUser := ctx.UserContext()

	ctxWithDeadline, cancel := context.WithDeadline(ctxUser, time.Now().Add(s.cfg.AppWriteTimeout))

	defer cancel()

	go func(chErr chan error, chWell chan struct{}) {
		if err := s.Service.DeleteSlideByID(ctxWithDeadline, param.ID, chWell); err != nil {
			s.log.Errorw("deleteSlide", "deleteSlide error", err.Error())

			chErr <- err

			close(chErr)
		}
	}(chErr, chWell)

	select {
	case <-chWell:
		return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success"))
	case err := <-chErr:
		if errors.Is(err, structs.ErrNoRowAffected) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	case <-ctxWithDeadline.Done():
		return ctx.Status(fiber.StatusRequestTimeout).JSON(structs.SetResponse(fiber.StatusRequestTimeout, fiber.ErrRequestTimeout.Message))
	}
}
