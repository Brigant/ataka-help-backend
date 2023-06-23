package api

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PartnerService interface {
	ReturnPartners(context.Context, structs.PartnerQueryParameters) ([]structs.Partner, int, error)
	SavePartner(context.Context, *multipart.Form, chan struct{}) error
	DeletePartnerByID(context.Context, string, chan struct{}) error
}

type Partner struct {
	Service PartnerService
	log     *logger.Logger
}

func NewParnerHandler(service PartnerService, log *logger.Logger) Partner {
	return Partner{
		Service: service,
		log:     log,
	}
}

func (p Partner) getPartners(ctx *fiber.Ctx) error {
	chErr := make(chan error)

	chWell := make(chan structs.PartnerResponse)

	params := structs.PartnerQueryParameters{
		Limit:  defaultLimit,
		Offset: defaultOffset,
	}

	if err := ctx.QueryParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	ctxUser := ctx.UserContext()

	ctxWithDeadline, cancel := context.WithDeadline(ctxUser, time.Now().Add(1*time.Second))

	defer cancel()

	go func(chErr chan error, chWell chan structs.PartnerResponse) {
		partners, total, err := p.Service.ReturnPartners(ctxWithDeadline, params)
		if err != nil && !errors.Is(err, structs.ErrNotFound) {
			p.log.Errorw("getPartner", "getPartner error", err.Error())

			chErr <- err

			close(chErr)
		}

		response := structs.PartnerResponse{
			Code:     fiber.StatusOK,
			Total:    total,
			Partners: partners,
		}

		chWell <- response

		close(chWell)

	}(chErr, chWell)

	/*
		sync.WaitGroup had not added because select{} blocks main goroutine any way.
	*/

	select {
	case response := <-chWell:
		return ctx.Status(fiber.StatusOK).JSON(response)
	case err := <-chErr:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	case <-ctxWithDeadline.Done():
		return ctx.Status(fiber.StatusRequestTimeout).JSON(structs.SetResponse(fiber.StatusRequestTimeout, fiber.ErrRequestTimeout.Message))
	}
}

func (p Partner) createPartner(ctx *fiber.Ctx) error {
	allowedFileExtentions := []string{"jpg", "jpeg", "webp", "png"}

	const (
		minAlt = 10
		maxAlt = 30
	)

	chErr := make(chan error)

	chWell := make(chan struct{})

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	alt := form.Value["alt"]

	if alt == nil || symbolsCounter(form.Value["alt"][0]) < minAlt || symbolsCounter(form.Value["alt"][0]) > maxAlt {
		p.log.Debugw("createPartner", "form.Value alt", "alt is blank or out of limits")

		return fiber.NewError(fiber.StatusBadRequest, "alt is blank or out of limits")
	}

	file := form.File["logo"]

	if file != nil {
		p.log.Debugw("createPartner", "file-name", file[0].Filename, "file-size", file[0].Size)
	}

	if file == nil || !isAllowedFileExtention(allowedFileExtentions, file[0].Filename) {
		p.log.Debugw("createPartner", "form.File", err.Error())

		return fiber.NewError(fiber.StatusBadRequest, "logo is absent")
	}

	size := file[0].Size

	if size > maxFileSize {
		p.log.Debugw("createPartner", "form.File", "file too large")

		return fiber.NewError(fiber.StatusBadRequest, "file too large")
	}

	ctxUser := ctx.UserContext()

	ctxWithDeadline, cancel := context.WithDeadline(ctxUser, time.Now().Add(1*time.Second))

	defer cancel()

	go func(chErr chan error, chWell chan struct{}) {
		if err := p.Service.SavePartner(ctxWithDeadline, form, chWell); err != nil {
			p.log.Errorw("createPartner", "createPartner error", err.Error())

			chErr <- err

			close(chErr)
		}
	}(chErr, chWell)

	/*
		sync.WaitGroup had not added because select{} blocks main goroutine any way.
	*/

	select {
	case _ = <-chWell:
		return ctx.JSON(structs.SetResponse(http.StatusOK, "success"))
	case err := <-chErr:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	case <-ctxWithDeadline.Done():
		return ctx.Status(fiber.StatusRequestTimeout).JSON(structs.SetResponse(fiber.StatusRequestTimeout, fiber.ErrRequestTimeout.Message))
	}
}

func (p Partner) deletePartner(ctx *fiber.Ctx) error {
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

	ctxWithDeadline, cancel := context.WithDeadline(ctxUser, time.Now().Add(1*time.Second))

	defer cancel()

	go func(chErr chan error, chWell chan struct{}) {
		if err := p.Service.DeletePartnerByID(ctxWithDeadline, param.ID, chWell); err != nil {
			p.log.Errorw("deletePartner", "deletePartner error", err.Error())

			chErr <- err

			close(chErr)
		}
	}(chErr, chWell)

	/*
		sync.WaitGroup had not added because select{} blocks main goroutine any way.
	*/

	select {
	case _ = <-chWell:
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
