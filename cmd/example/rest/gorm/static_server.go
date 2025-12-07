package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"strings"
)

func ServeStaticSPA(r *gin.Engine, urlPrefix, distDir string) {
	if urlPrefix == "" {
		urlPrefix = ""
	} else {
		if !strings.HasPrefix(urlPrefix, "/") {
			urlPrefix = "/" + urlPrefix
		}
		urlPrefix = strings.TrimRight(urlPrefix, "/")
	}

	assetsRoute := urlPrefix + "/assets"
	r.Static(assetsRoute, filepath.Join(distDir, "assets"))

	indexPath := filepath.Join(distDir, "index.html")
	r.GET(urlPrefix+"/", func(c *gin.Context) {
		c.File(indexPath)
	})
	r.GET(urlPrefix+"/index.html", func(c *gin.Context) {
		c.File(indexPath)
	})

	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			NotFoundHandler(c)
			return
		}

		p := c.Request.URL.Path
		if urlPrefix != "" && strings.HasPrefix(p, urlPrefix) {
			p = strings.TrimPrefix(p, urlPrefix)
		}
		last := filepath.Base(p)
		if strings.Contains(last, ".") {
			NotFoundHandler(c)
			return
		}

		c.File(indexPath)
	})
}
