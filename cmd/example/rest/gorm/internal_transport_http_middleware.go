package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info().Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Msg("request start")
		c.Next()
		log.Info().Int("status", c.Writer.Status()).Msg("request done")
	}
}

func JSONError(c *gin.Context, status int, err string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": err,
	})
}

func NotFoundHandler(c *gin.Context) {
	JSONError(c, http.StatusNotFound, "resource not found")
}
