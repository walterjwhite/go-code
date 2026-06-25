package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	chatlib "github.com/walterjwhite/go-code/lib/chat"
)

var (
	contactLimiter = chatlib.NewIPRateLimiter(chatlib.RateLimitConfig{
		Messages: 5,
		Window:   10 * time.Minute,
	})
	tokenLimiter = chatlib.NewIPRateLimiter(chatlib.RateLimitConfig{
		Messages: 30,
		Window:   time.Minute,
	})
)

func serve() *http.Server {
	router := gin.Default()

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatal().Err(err).Msg("failed to configure trusted proxies")
	}

	router.Use(securityHeadersMiddleware)
	router.Use(gin.CustomRecovery(recoveryHandler), gin.Logger())

	contactLimiter.Start(context.Background(), 10*time.Minute)
	tokenLimiter.Start(context.Background(), time.Minute)

	httpRoutes := router.Group("/")
	httpRoutes.Use(maxBytesMiddleware(64 * 1024))
	{
		httpRoutes.POST("/api/contact", ipRateLimitMiddleware(contactLimiter), onContactRequest)
		httpRoutes.POST("/api/token", ipRateLimitMiddleware(tokenLimiter), onTokenEvent)
	}

	router.GET("/ws/chat", onChatWSRequest)

	if chatClient != nil {
		chatHub.start(context.Background(), chatClient)
	}

	addr := fmt.Sprintf("%s:%d", *hostFlag, *portFlag)

	return &http.Server{
		Addr:    addr,
		Handler: router,

		ReadTimeout: 10 * time.Second,

		WriteTimeout: 30 * time.Second,

		IdleTimeout: 120 * time.Second,
	}
}

func maxBytesMiddleware(n int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, n)
		c.Next()
	}
}

func ipRateLimitMiddleware(limiter *chatlib.IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(c.ClientIP()) {
			JSONError(c, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		c.Next()
	}
}

func securityHeadersMiddleware(c *gin.Context) {
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy",
		"default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
	c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	c.Next()
}

func recoveryHandler(c *gin.Context, recovered any) {
	log.Error().Interface("panic", recovered).Msg("internal server error")
	JSONError(c, http.StatusInternalServerError, "internal server error")
}

func JSONError(c *gin.Context, status int, err string) {
	c.AbortWithStatusJSON(status, gin.H{"error": err})
}
