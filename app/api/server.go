package api

import (
	"context"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	HTTPServer *fiber.App
}

func NewServer(cfg config.Config) *Server {
	server := new(Server)
	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
	}

	server.HTTPServer = fiber.New(fconfig)

	return server
}

type Serv struct {
	HTTPServer *fiber.App
	serices    ServiceInterfaces
}

func NewServ(cfg config.Config, services ServiceInterfaces) *Serv {
	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
	}

	server := &Serv{
		HTTPServer: fiber.New(fconfig),
		serices:    services,
	}

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.ShutdownWithContext(ctx)
}
