package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	srv *http.Server
}

func NewServer(router *gin.Engine, cfg *Config) *Server {
	addr := ":" + cfg.AppPort
	s := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.ReadTimeout,
		IdleTimeout:  60 * time.Second,
	}
	return &Server{srv: s}
}

func (s *Server) Start() error {
	log.Info().Str("addr", s.srv.Addr).Msg("starting HTTP server")
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("shutting down HTTP server")
	return s.srv.Shutdown(ctx)
}
