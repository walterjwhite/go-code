package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"net/http"

	"time"
)

func serve() *http.Server {
	router := gin.Default()

	router.Use(securityHeadersMiddleware)

	router.Use(gin.CustomRecovery(recoveryHandler), gin.Logger())

	router.POST("/api/contact", onContactRequest)

	addr := fmt.Sprintf("%s:%d", *hostFlag, *portFlag)

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func securityHeadersMiddleware(c *gin.Context) {
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
	c.Next()
}

func recoveryHandler(c *gin.Context, recovered any) {
	log.Error().Interface("panic", recovered).Msg("internal server error")

	JSONError(c, http.StatusInternalServerError, "internal server error")
}

func JSONError(c *gin.Context, status int, err string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": err,
	})
}
