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

const bodyLimit = 5 * 1024 * 1024

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

	app.Get("/cards", timeout.NewWithContext(h.Card.getCards, cfg.Server.AppReadTimeout))
	app.Post("/cards", identifyUser, timeout.NewWithContext(h.Card.createCard, cfg.Server.AppWriteTimeout))
	app.Get("/cards/:id", timeout.NewWithContext(h.Card.findCard, cfg.Server.AppReadTimeout))
	app.Delete("/cards/:id", identifyUser, timeout.NewWithContext(h.Card.deleteCard, cfg.Server.AppWriteTimeout))

	app.Get("/partners", h.Partner.getPartners)
	app.Post("/partners", h.Partner.createPartner)
	app.Delete("/partners/:id", h.Partner.deletePartner)

	app.Get("/slider", h.Slider.getSlider)
	app.Post("/slider", identifyUser, h.Slider.createSlider)
	app.Delete("/slider/:id", identifyUser, h.Slider.deleteSlide)

	app.Put("/contacts", identifyUser, timeout.NewWithContext(h.Contact.edit, cfg.Server.AppReadTimeout))
	app.Get("/contacts", timeout.NewWithContext(h.Contact.get, cfg.Server.AppReadTimeout))

	app.Get("/reports", h.Report.getReports)
	app.Put("/reports", identifyUser, h.Report.updateReport)
	app.Delete("/reports", identifyUser, h.Report.deleteReport)

	app.Post("/feedback", timeout.NewWithContext(h.Feedback.sendFedback, cfg.Server.AppWriteTimeout))

	app.Post("/auth/login", timeout.NewWithContext(h.Auth.login, cfg.Server.AppWriteTimeout))
	app.Post("/auth/logout", identifyUser, timeout.NewWithContext(h.Auth.logout, cfg.Server.AppWriteTimeout))
	app.Post("/auth/refresh", h.Auth.refresh)
	app.Post("/auth/change", identifyUser, timeout.NewWithContext(h.Auth.change, cfg.Server.AppReadTimeout))
}

func corsConfig() cors.Config {
	return cors.Config{
		// AllowOrigins: `https://ataka-help.vercel.app, http://localhost,  http://localhost:7000`,
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}
}
