package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const bodyLimit = 2 * 1024 * 1024

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

	server.initRoutes(server.HTTPServer, handler)

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return fmt.Errorf("shutdown error: %w", s.HTTPServer.ShutdownWithContext(ctx))
}

func (s Server) initRoutes(app *fiber.App, h Handler) {
	app.Static("/static", "./static")

	app.Get("/cards", h.Card.getCards)
	app.Post("/cards", h.Card.createCard)
	app.Get("/cards/:id", h.Card.findCard)
	app.Delete("cards/:id", h.Card.deleteCard)

	app.Get("/partners", h.Partner.get)

	app.Get("/slider", h.Slider.getSlider)
	app.Post("/slider", h.Slider.createSlider)

	app.Put("/contacts", h.Contact.edit)
	app.Get("/contacts", h.Contact.get)

	app.Get("/reports", h.Report.getReports)
	app.Put("/reports", h.Report.updateReport)
	app.Delete("/reports", h.Report.deleteReport)

	app.Post("/feedback", h.Feedback.sendFedback)
}

func corsConfig() cors.Config {
	return cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}
}
