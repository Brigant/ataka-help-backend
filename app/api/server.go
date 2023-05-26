package api

import (
	"context"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	HTTPServer *fiber.App
}

func NewServer(cfg config.Config, handler Handler) *Server {
	server := new(Server)
	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
		BodyLimit:    5 * 1024 * 1024,
	}

	server.HTTPServer = fiber.New(fconfig)

	server.initRoutes(server.HTTPServer, handler)

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.ShutdownWithContext(ctx)
}

func (s Server) initRoutes(app *fiber.App, h Handler) {
	app.Get("/", h.Card.getCards)
	app.Static("/", "../static")
	app.Post("/", h.Card.createCard)
	app.Get("/partners", h.Partner.Get)
}
