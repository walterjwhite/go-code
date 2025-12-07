package request_filter

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"net/http"
	"strings"
	"sync"
	"time"
)

type Conf struct {
	Tokens []string

	AlternativePath string

	TTL        time.Duration
	visitorsMu sync.RWMutex
	visitors   map[string]time.Time

	timers map[string]*time.Timer
	AllowPrefixes []string
}

func (i *Conf) Handler() gin.HandlerFunc {
	if i.visitors == nil {
		i.visitors = make(map[string]time.Time)
	}
	if i.timers == nil {
		i.timers = make(map[string]*time.Timer)
	}
	if i.AlternativePath == "" {
		i.AlternativePath = "./static/alternative.html"
		if !fileExists(i.AlternativePath) {
			log.Warn().Msgf("linkFilter: alternative content file %s not found", i.AlternativePath)
			i.AlternativePath = ""
		}
	}

	if len(i.AllowPrefixes) == 0 {
		i.AllowPrefixes = []string{"/static/error/", "/static/alternative.html", "/favicon.ico", "/favicon-", "/site.webmanifest", "/android-chrome-", "/apple-touch-icon.png", "/robots.txt", "/app.js", "/style.css"}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, p := range i.AllowPrefixes {
			if p == "" {
				continue
			}
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}
		if c.Request.Method != http.MethodGet {
			log.Warn().Msgf("linkFilter: non-GET request %s %s", c.Request.Method, c.Request.URL.Path)
			c.Next()
			return
		}

		ip := c.ClientIP()
		if i.visitorSeen(ip) {
			log.Debug().Msgf("visitorSeen: allowing IP %s", ip)
			c.Next()
			return
		}

		if val := strings.TrimSpace(c.Query("token")); val != "" {
			for _, t := range i.Tokens {
				if t == val {
					log.Warn().Msgf("token bypass: allowing IP %s", ip)
					i.markVisitorSeen(ip)
					c.Next()
					return
				}
			}
		}

		if _, err := c.Writer.WriteString(""); err == nil {
			if len(i.AlternativePath) > 0 {
				log.Warn().Msgf("serving alternative content to IP %s", ip)
				c.File(i.AlternativePath)
			} else {
				log.Warn().Msgf("serving alternative content to IP - path does not exist: %s", ip)
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<html><body><h1>Welcome!</h1><p>Please use the special link to view the site.</p></body></html>"))
			}
		}
		c.Abort()
	}
}
