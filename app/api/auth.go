package api

import (
	"fmt"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AutService interface{}

type AuthHandler struct {
	Service AutService
	log     *logger.Logger
}

func NewAuthHandler(service AutService, log *logger.Logger) AuthHandler {
	return AuthHandler{
		Service: service,
		log:     log,
	}
}

func (h AuthHandler) login(ctx *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "john"
	cookie.Value = uuid.NewString()
	// cookie.Domain = os.Getenv("SRV_HOST")
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	// cookie.Secure = true

	// Set cookie
	ctx.Cookie(cookie)

	ctx.Context()
	fmt.Println()

	return ctx.Status(fiber.StatusOK).JSON("http.SameSiteLaxMode")
}

func (h AuthHandler) refresh(ctx *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "john"
	cookie.Value = uuid.NewString()
	// cookie.Domain = os.Getenv("SRV_HOST")
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	// cookie.Secure = true

	// Set cookie
	ctx.Cookie(cookie)
	return ctx.Status(fiber.StatusOK).JSON("http.SameSiteLaxMode")
}
