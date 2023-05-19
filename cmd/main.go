package main

import (
	"log"

	"github.com/baza-trainee/ataka-help-backend/app/api"
	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/repository/pg"
	"github.com/baza-trainee/ataka-help-backend/app/services"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Print(err.Error())
	}

	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Print(err.Error())
	}

	db, err := pg.NewPostgresDB(cfg)
	if err != nil {
		log.Print(err.Error())
	}

	repo := pg.NewRepository(db)

	service := services.NewService(repo)

	handler := api.NewHandler(service, logger)

	server := api.NewServer(cfg)

	handler.InitRoutes(server.HttpServer)

	if err := server.HttpServer.Listen(cfg.Server.AppAddress); err != nil {
		log.Println(err.Error())
	}
}
