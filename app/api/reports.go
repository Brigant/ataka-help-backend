package api

import (
	"errors"
	"mime/multipart"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

type ReportService interface {
	ReturnReport() (string, error)
	ChangeReport(*multipart.Form) error
	DeleteReport() error
}

type ReportHandler struct {
	Service ReportService
	log     *logger.Logger
}

func NewReportHandler(service ReportService, log *logger.Logger) ReportHandler {
	return ReportHandler{
		Service: service,
		log:     log,
	}
}

func (h ReportHandler) getReports(ctx *fiber.Ctx) error {
	report, err := h.Service.ReturnReport()
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	response := structs.ReportResponse{
		Code: fiber.StatusOK,
		File: report,
	}

	return ctx.Status(fiber.StatusOK).JSON(response) // nolint
}

func (h ReportHandler) updateReport(ctx *fiber.Ctx) error {
	const minimalCountItems = 1
	allowedType := []string{"application/pdf"}
	allowedExtention := []string{"pdf"}

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if len(form.File["thumb"]) < minimalCountItems {
		h.log.Debugw("updateReport", "form.File", "no repport was attached")

		return fiber.NewError(fiber.StatusBadRequest, "no repport was attached")
	}

	fileHeader := form.File["thumb"][0]

	if fileHeader == nil || fileHeader.Size > fileLimit || !isAllowedContentType(allowedType, fileHeader.Header["Content-Type"][0]) {
		h.log.Debugw("updateReport", "form.File", "required file not bigger then 2 Mb and in pdf format")

		return fiber.NewError(fiber.StatusBadRequest, "required file not bigger then 2 Mb and in pdf format")
	}

	if !isAllowedFileExtention(allowedExtention, fileHeader.Filename) {
		h.log.Debugw("updateReport", "isAllowedFileExtention", "required file in pdf format")

		return fiber.NewError(fiber.StatusBadRequest, "required file in pdf format")
	}

	if err := h.Service.ChangeReport(form); err != nil {
		h.log.Errorf("ChangeReport", "error", err.Error())

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(structs.SetResponse(fiber.StatusCreated, "success")) // nolint
}

func (h ReportHandler) deleteReport(ctx *fiber.Ctx) error {
	if err := h.Service.DeleteReport(); err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(structs.SetResponse(fiber.StatusOK, "success")) // nolint
}
