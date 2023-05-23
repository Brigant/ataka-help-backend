package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baza-trainee/ataka-help-backend/app/api"
	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/logger"
	"github.com/baza-trainee/ataka-help-backend/app/pg"
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

	repo, err := pg.NewRepository(cfg)
	if err != nil {
		logger.Errorw("New Repository", "error", err.Error())
	}

	service := services.NewService(repo)

	handler := api.NewHandler(service, logger)

	server := api.NewServer(cfg)

	handler.InitRoutes(server.HTTPServer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.HTTPServer.Listen(cfg.Server.AppAddress); err != nil {
			log.Println(err.Error())
		}
	}()

	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error closing server: %v", err)
	}

	if err := repo.Close(); err != nil {
		logger.Errorf("error occured on db connection close: %s", err.Error())
	}

	log.Println("server stopped")
}
