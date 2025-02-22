package main

import (
	"log"
	"net/http"

	"github.com/dqx0/hitandblow/internal/server"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/dqx0/hitandblow/docs" // generated docs package
)

// @title Hit and Blow API
// @version 1.0.0
// @description Hit and Blow GameのAPIドキュメント
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	// Swagger UI を /swagger/index.html で公開
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gameServer := server.NewGameServer()

	// ヘルスチェック用エンドポイント
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// ゲーム開始エンドポイント
	r.POST("/games", gameServer.StartNewGame)

	// 判定エンドポイント
	r.POST("/games/:gameId/guess", gameServer.MakeGuess)

	log.Println("サーバーを開始します: http://localhost:8080")
	r.Run(":8080")
}
