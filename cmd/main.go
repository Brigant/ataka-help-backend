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
	const timoutLimit = 5

	cfg, err := config.InitConfig()
	if err != nil {
		log.Print(err.Error())
	}

	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Print(err.Error())
	}

	defer logger.Flush()

	repo, err := pg.NewRepository(cfg)
	if err != nil {
		logger.Errorw("New Repository", "error", err.Error())
	}

	service := services.NewService(repo, cfg.Smtp)

	handler := api.NewHandler(service, logger)

	server := api.NewServer(cfg, handler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.HTTPServer.Listen(cfg.Server.AppAddress); err != nil {
			logger.Errorw("Start and Listen", "error", err.Error())
		}
	}()

	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timoutLimit*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorw("Shutdown server", "error", err.Error())
	}

	if err := repo.Close(); err != nil {
		logger.Errorw("Close repository", "error", err.Error())
	}

	log.Println("server stopped")
}
