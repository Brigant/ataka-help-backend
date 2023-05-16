package rest

import (
	"fmt"

	"github.com/baza-trainee/ataka-help-backend/app/repositories/pg"
	"github.com/baza-trainee/ataka-help-backend/app/rest/api"
	"github.com/baza-trainee/ataka-help-backend/app/services"
	"github.com/baza-trainee/ataka-help-backend/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	httpServer *fiber.App
}

func NewServer(cfg config.Config) *Server {
	server := new(Server)
	config := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
	}

	server.httpServer = fiber.New(config)

	return server
}

// Compose and init all dependensies and start server
func SetupAndRun() error {
	cfg, err := config.InitConfig()
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	db, err := pg.NewPostgresDB(cfg)
	if err != nil {
		return fmt.Errorf("error while creating connection to database: %w", err)
	}

	storage := pg.NewRepository(db)

	service := services.New(
		services.Deps{
			CardsStorage: storage.CardsDB,
		})

	handler := api.NewHandler(
		api.Deps{
			CardsService: service.CardsService,
		})

	app := handler.InitRoutes(cfg)

	app.Listen(cfg.Server.Port)
	return nil
}
