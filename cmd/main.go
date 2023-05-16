package main

import (
	"log"
	"os"

	"github.com/baza-trainee/ataka-help-backend/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	env, err := os.Open(".env")
	if err != nil {
		log.Println(err.Error())
	}

	defer func() {
		if err = env.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	config, err := config.InitConfig()
	if err != nil {
		log.Println(err.Error())
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	_ = config

	if err := app.Listen(config.Server.AppAddress); err != nil {
		log.Println(err.Error())
	}
}
