package main

import (
	"net/http"
	"nh-be/config"

	"github.com/gin-gonic/gin"
)

func main() {

  router := gin.Default()

  db := config.ConnectDatabase()
  defer db.Close()

  router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  router.Run()
}