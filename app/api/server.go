package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/api/midlware"
	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

const bodyLimit = 10 * 1024 * 1024

type Server struct {
	HTTPServer *fiber.App
}

func NewServer(cfg config.Config, handler Handler) *Server {
	server := new(Server)

	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
		BodyLimit:    bodyLimit,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			ctx.Status(code)

			if err := ctx.JSON(err); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			return nil
		},
	}

	server.HTTPServer = fiber.New(fconfig)

	server.HTTPServer.Use(cors.New(corsConfig()))

	server.HTTPServer.Use(logger.New())

	server.HTTPServer.Use(recover.New())

	server.initRoutes(server.HTTPServer, handler, cfg)

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return fmt.Errorf("shutdown error: %w", s.HTTPServer.ShutdownWithContext(ctx))
}

func (s Server) initRoutes(app *fiber.App, h Handler, cfg config.Config) {
	identifyUser := midlware.NewUserIdentity(cfg.Auth)

	app.Static("/static", "./static")

	api := app.Group(apiPrefixV1)
	{
		api.Get("/cards", timeout.NewWithContext(h.Card.getCards, cfg.Server.AppReadTimeout))
		api.Post("/cards", identifyUser, timeout.NewWithContext(h.Card.createCard, cfg.Server.AppWriteTimeout))
		api.Get("/cards/:id", timeout.NewWithContext(h.Card.findCard, cfg.Server.AppReadTimeout))
		api.Delete("/cards/:id", identifyUser, timeout.NewWithContext(h.Card.deleteCard, cfg.Server.AppWriteTimeout))

		api.Get("/partners", h.Partner.getPartners)
		api.Post("/partners", identifyUser, h.Partner.createPartner)
		api.Delete("/partners/:id", identifyUser, h.Partner.deletePartner)

		api.Get("/slider", h.Slider.getSlider)
		api.Post("/slider", identifyUser, h.Slider.createSlide)
		api.Delete("/slider/:id", identifyUser, h.Slider.deleteSlide)

		api.Put("/contacts", identifyUser, timeout.NewWithContext(h.Contact.edit, cfg.Server.AppReadTimeout))
		api.Get("/contacts", timeout.NewWithContext(h.Contact.get, cfg.Server.AppReadTimeout))

		api.Get("/reports", h.Report.getReports)
		api.Put("/reports", identifyUser, h.Report.updateReport)
		api.Delete("/reports", identifyUser, h.Report.deleteReport)

		api.Post("/feedback", timeout.NewWithContext(h.Feedback.sendFedback, cfg.Server.AppWriteTimeout))

		api.Post("/auth/login", timeout.NewWithContext(h.Auth.login, cfg.Server.AppWriteTimeout))
		api.Post("/auth/logout", identifyUser, timeout.NewWithContext(h.Auth.logout, cfg.Server.AppWriteTimeout))
		api.Post("/auth/refresh", h.Auth.refresh)
		api.Post("/auth/change", identifyUser, timeout.NewWithContext(h.Auth.change, cfg.Server.AppReadTimeout))
	}
}

func corsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     `https://ataka-help.tech, http://localhost:3000`,
		AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Credentials",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}
}
