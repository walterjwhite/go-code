package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func init() {
	Configure()

	log.Debug().Msg("initializing database")
}

func cleanup(database *gorm.DB, server *Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	Close(database)
	log.Info().Msg("shutdown complete")
}

func main() {
	defer OnPanic()
	cfg, err := Load()
	Error(err)

	database, err := NewSQLite(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to sqlite db")
	}

	if err := database.AutoMigrate(&User{}); err != nil {
		log.Fatal().Err(err).Msg("auto-migrate failed")
	}

	userRepo := NewGormUserRepository(database)
	userSvc := NewUserService(userRepo)
	handler := NewHandler(userSvc)

	server := NewServer(handler.Router(cfg.DatabaseURL), cfg)
	defer cleanup(database, server)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	log.Info().Msgf("server started on :%s", cfg.AppPort)

	Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	CloseSQLite(database)
	log.Info().Msg("shutdown complete")
}
