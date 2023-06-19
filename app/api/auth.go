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
}

type AuthHandler struct {
	Service AutService
	log     *logger.Logger
	Auth    config.AuthConfig
}

func NewAuthHandler(service AutService, log *logger.Logger, cfg config.AuthConfig) AuthHandler {
	return AuthHandler{
		Service: service,
		log:     log,
		Auth:    cfg,
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

	tokenPair, err := h.Service.GetTokenPair(ctx.UserContext(), identity, h.Auth)
	if err != nil {
		if errors.Is(err, structs.ErrTimeout) {
			return fiber.NewError(fiber.StatusRequestTimeout, err.Error())
		}

		if errors.Is(err, structs.ErrNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessCookie := new(fiber.Cookie)
	accessCookie.Name = structs.AccessCookieName
	accessCookie.Value = tokenPair.AccessToken
	accessCookie.Expires = tokenPair.AccessExpire
	accessCookie.HTTPOnly = true
	// cookie.Secure = true

	ctx.Cookie(accessCookie)

	refreshCookie := new(fiber.Cookie)
	refreshCookie.Name = structs.RefreshCookieName
	refreshCookie.Value = tokenPair.RefreshToken
	refreshCookie.Expires = tokenPair.RefresgExpire
	refreshCookie.HTTPOnly = true
	// cookie.Secure = true

	ctx.Cookie(refreshCookie)

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "logined"))
}

func (h AuthHandler) refresh(ctx *fiber.Ctx) error {
	refreshString := ctx.Cookies(structs.RefreshCookieName)

	if refreshString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "empty cookie")
	}

	userID, err := midlware.ParseToken(refreshString, h.Auth.SigningKey)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	tokenPair, err := h.Service.Refresh(refreshString, userID, h.Auth)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessCookie := new(fiber.Cookie)
	accessCookie.Name = structs.AccessCookieName
	accessCookie.Value = tokenPair.AccessToken
	accessCookie.Expires = tokenPair.AccessExpire
	accessCookie.HTTPOnly = true
	// cookie.Secure = true

	ctx.Cookie(accessCookie)

	refreshCookie := new(fiber.Cookie)
	refreshCookie.Name = structs.RefreshCookieName
	refreshCookie.Value = tokenPair.RefreshToken
	refreshCookie.Expires = tokenPair.RefresgExpire
	refreshCookie.HTTPOnly = true
	// cookie.Secure = true

	ctx.Cookie(refreshCookie)

	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "refreshed"))
}

func (h AuthHandler) logout(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID")

	accesCookie := new(fiber.Cookie)
	accesCookie.Name = structs.AccessCookieName
	accesCookie.Value = ""
	accesCookie.Expires = time.Now().Add(-1)
	accesCookie.HTTPOnly = true
	// cookie.Secure = true

	refreshCookie := new(fiber.Cookie)
	refreshCookie.Name = structs.RefreshCookieName
	refreshCookie.Value = ""
	refreshCookie.Expires = time.Now().Add(-1)
	refreshCookie.HTTPOnly = true
	// cookie.Secure = true

	h.Service.CleanSession(userID.(string))

	ctx.Cookie(accesCookie)
	ctx.Cookie(refreshCookie)
	return ctx.Status(fiber.StatusOK).JSON(structs.SetResponse(fiber.StatusOK, "success"))
}
