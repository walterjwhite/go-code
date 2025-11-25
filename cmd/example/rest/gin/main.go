package main

import (
	"github.com/gin-gonic/gin"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	logging.Panic(router.Run("localhost:8080"))
}
