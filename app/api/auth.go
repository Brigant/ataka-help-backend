package api

import (
	"context"
	"errors"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

type AutService interface {
	GetTokenPair(context.Context, structs.IdentityData, config.AuthConfig) (string, time.Time, error)
	CleanSession(userID string)
	Refresh(refreshString string, cfg config.AuthConfig) (string, time.Time, error)
}

type AuthHandler struct {
	Service    AutService
	log        *logger.Logger
	Auth       config.AuthConfig
	CookieName string
}

func NewAuthHandler(service AutService, log *logger.Logger, cfg config.AuthConfig) AuthHandler {
	return AuthHandler{
		Service:    service,
		log:        log,
		Auth:       cfg,
		CookieName: "TokenPair",
	}
}

func (h AuthHandler) login(ctx *fiber.Ctx) error {
	identity := structs.IdentityData{}

	if err := ctx.BodyParser(&identity); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, structs.ErrWrongBody.Error())
	}

	if identity.Login == "" || identity.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, structs.ErrWrongBody.Error())
	}

	tokenPairJSON, expires, err := h.Service.GetTokenPair(ctx.UserContext(), identity, h.Auth)
	if err != nil {
		if errors.Is(err, structs.ErrTimeout) {
			return fiber.NewError(fiber.StatusRequestTimeout, err.Error())
		}

		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	cookie := new(fiber.Cookie)
	cookie.Name = h.CookieName
	cookie.Value = tokenPairJSON
	cookie.Expires = expires
	cookie.HTTPOnly = true
	// cookie.Secure = true

	ctx.Cookie(cookie)

	return ctx.Status(fiber.StatusOK).JSON(identity)
}

func (h AuthHandler) refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Locals("refreshString")
	if refreshToken == nil || refreshToken == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "cann't find refreshString in context")
	}

	tokenPairJSON, expires, err := h.Service.Refresh(refreshToken.(string), h.Auth)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	cookie := new(fiber.Cookie)
	cookie.Name = h.CookieName
	cookie.Value = tokenPairJSON
	cookie.Expires = expires
	cookie.HTTPOnly = true
	// cookie.Secure = true

	ctx.Cookie(cookie)

	return ctx.Status(fiber.StatusOK).JSON("refreshed")
}

func (h AuthHandler) logout(ctx *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = h.CookieName
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1)
	cookie.HTTPOnly = true
	// cookie.Secure = true

	// TODO: remove refresh token from memory

	ctx.Cookie(cookie)
	return ctx.Status(fiber.StatusOK).JSON("ok")
}
