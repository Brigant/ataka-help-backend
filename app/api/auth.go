package api

import (
	"context"
	"errors"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/api/midlware"
	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

type AutService interface {
	GetTokenPair(context.Context, structs.IdentityData, config.AuthConfig) (structs.TokenPair, error)
	CleanSession(string)
	Refresh(string, string, config.AuthConfig) (structs.TokenPair, error)
	ChangePassword(context.Context, string, structs.PasswordsContainer, config.AuthConfig) error
}

type AuthHandler struct {
	Service    AutService
	log        *logger.Logger
	AuthConfig config.AuthConfig
}

func NewAuthHandler(service AutService, log *logger.Logger, cfg config.AuthConfig) AuthHandler {
	return AuthHandler{
		Service:    service,
		log:        log,
		AuthConfig: cfg,
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

	tokenPair, err := h.Service.GetTokenPair(ctx.UserContext(), identity, h.AuthConfig)
	if err != nil {
		if errors.Is(err, structs.ErrTimeout) {
			return fiber.NewError(fiber.StatusRequestTimeout, err.Error())
		}

		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessCookie := newCookie(
		structs.AccessCookieName,
		tokenPair.AccessToken,
		tokenPair.AccessExpire,
	)

	ctx.Cookie(accessCookie)

	refreshCookie := newCookie(
		structs.RefreshCookieName,
		tokenPair.RefreshToken,
		tokenPair.RefresgExpire,
	)

	refreshCookie.Path = "/auth/refresh"

	ctx.Cookie(refreshCookie)

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "logined"))
}

func (h AuthHandler) refresh(ctx *fiber.Ctx) error {
	refreshString := ctx.Cookies(structs.RefreshCookieName)

	if refreshString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "empty cookie")
	}

	userID, err := midlware.ParseToken(refreshString, h.AuthConfig.SigningKey)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tokenPair, err := h.Service.Refresh(refreshString, userID, h.AuthConfig)
	if err != nil {
		if errors.Is(err, structs.ErrNoSession) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessCookie := newCookie(
		structs.AccessCookieName,
		tokenPair.AccessToken,
		tokenPair.AccessExpire,
	)

	refreshCookie := newCookie(
		structs.RefreshCookieName,
		tokenPair.RefreshToken,
		tokenPair.RefresgExpire,
	)

	refreshCookie.Path = "/auth/refresh"

	ctx.Cookie(accessCookie)
	ctx.Cookie(refreshCookie)

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "refreshed"))
}

func (h AuthHandler) logout(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID")

	h.Service.CleanSession(userID.(string))

	accessCookie := newCookie(
		structs.AccessCookieName,
		"",
		time.Now().Add(-1),
	)

	refreshCookie := newCookie(
		structs.RefreshCookieName,
		"",
		time.Now().Add(-1),
	)
	refreshCookie.Path = "/auth/refresh"
	ctx.Cookie(accessCookie)
	ctx.Cookie(refreshCookie)
	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success"))
}

func (h AuthHandler) change(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID")

	passContainer := structs.PasswordsContainer{}

	if err := ctx.BodyParser(&passContainer); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := passContainer.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.Service.ChangePassword(ctx.UserContext(), userID.(string), passContainer, h.AuthConfig); err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "wrong current password")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success"))
}

func newCookie(name, value string, expire time.Time) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		Secure:   false,
		HTTPOnly: true,
		SameSite: "none",
	}
}
