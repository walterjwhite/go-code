package request_logger

import (
	"github.com/rs/zerolog/log"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Conf) Handler() gin.HandlerFunc {
	out := c.initRequestLogChannel()

	return func(c *gin.Context) {
		c.Next()

		rl := RequestLog{
			TS:         time.Now().UTC(),
			IP:         clientIP(c.Request.Header.Get("X-Forwarded-For"), c.ClientIP()),
			Method:     c.Request.Method,
			RequestURI: c.Request.RequestURI,
			UserAgent:  c.Request.UserAgent(),
			Status:     c.Writer.Status(),
		}

		select {
		case out <- rl:
		default:
			log.Warn().Msg("request log dropped; buffer full")
		}
	}
}

func clientIP(forwarded string, fallback string) string {
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				if net.ParseIP(p) != nil {
					return p
				}
			}
		}
	}
	host, _, err := net.SplitHostPort(fallback)
	if err == nil {
		return host
	}
	return fallback
}
