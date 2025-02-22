package main

import (
	"log"
	"net/http"

	"github.com/dqx0/hitandblow/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gameServer := server.NewGameServer()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/games", gameServer.StartNewGame)
	r.POST("/games/:gameId/guess", gameServer.MakeGuess)

	log.Println("サーバーを開始します: http://localhost:8080")
	r.Run(":8080")
}
