package api

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	httpServer *fiber.App
}

func NewServer() *Server {
	server := new(Server)
	config := fiber.Config{
		ReadTimeout:  60,
		WriteTimeout: 60,
		IdleTimeout:  60,
	}

	server.httpServer = fiber.New(config)

	return server
}
