package api

import (
	"mime/multipart"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

type ReportService interface {
	ReturnReport() string
	ChangeReport(*multipart.Form) error
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
	report := h.Service.ReturnReport()

	response := structs.ReportResponse{
		Code: fiber.StatusOK,
		File: report,
	}

	return ctx.Status(fiber.StatusOK).JSON(response) // nolint
}

func (h ReportHandler) updateReport(ctx *fiber.Ctx) error {
	allowedFileExtentions := []string{"pdf"}

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if len(form.File["report"]) < 1 {
		h.log.Debugw("updateReport", "form.File", "no repport was attached")

		return fiber.NewError(fiber.StatusBadRequest, "no repport was attached")
	}

	fileHeader := form.File["report"][0]

	if fileHeader == nil || fileHeader.Size > fileLimit || !isAllowedFileExtention(allowedFileExtentions, fileHeader.Filename) {
		h.log.Debugw("updateReport", "form.File", "required file not bigger then 5 Mb and in pdf format")

		return fiber.NewError(fiber.StatusBadRequest, "required file not bigger then 5 Mb and in pdf format")
	}

	if err := h.Service.ChangeReport(form); err != nil {
		h.log.Errorf("ChangeReport", "error", err.Error())

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(structs.SetResponse(fiber.StatusCreated, "success")) // nolint
}
