package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	port := "8080"

	// starts a new gin instance with no middle-ware
	r := gin.New()

	// define handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// listen and serve on defined port
	logrus.Debugf("Listening on port %s", port)
	r.Run(":" + port)

}
