package main

import (
	"github.com/gin-gonic/gin"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
	"sync"
)

func serve(wg *sync.WaitGroup) {
	defer wg.Done()

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.POST("/api/contact", contact)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "localhost:" + port
	logging.Error(router.Run(addr))
}
