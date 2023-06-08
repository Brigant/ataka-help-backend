package api

import (
	"context"
	"errors"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

type FeedbackService interface {
	PassFeedback(context.Context, structs.Feedback) error
}

type FeedbackHandler struct {
	Service FeedbackService
	log     *logger.Logger
}

func NewFeedbackHandler(service FeedbackService, log *logger.Logger) FeedbackHandler {
	return FeedbackHandler{
		Service: service,
		log:     log,
	}
}

func (h FeedbackHandler) sendFedback(ctx *fiber.Ctx) error {
	feedback := structs.Feedback{}
	if err := ctx.BodyParser(&feedback); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := feedback.Valiadate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.Service.PassFeedback(ctx.Context(), feedback); err != nil {
		if errors.Is(err, structs.ErrCheckCaptcha) {
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		}

		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success"))
}
