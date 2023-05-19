package api

import (
	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	HttpServer *fiber.App
}

func NewServer(cfg config.Config) *Server {
	server := new(Server)
	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
	}

	server.HttpServer = fiber.New(fconfig)

	return server
}
