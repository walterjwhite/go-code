package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TokenEvent struct {
	Token     string    `json:"token"`
	IP        string    `json:"ip"`
	Referrer  string    `json:"referrer,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func onTokenEvent(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Referrer string `json:"referrer"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	token := strings.TrimSpace(req.Token)
	if len(token) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token too long"})
		return
	}

	event := TokenEvent{
		Token:     token,
		IP:        c.ClientIP(),
		Referrer:  req.Referrer,
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now().UTC(),
	}

	log.Info().
		Str("token", event.Token).
		Str("ip", event.IP).
		Str("referrer", event.Referrer).
		Str("user_agent", event.UserAgent).
		Time("timestamp", event.Timestamp).
		Msg("resume token visit")

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
