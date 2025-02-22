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
	// @Summary ヘルスチェック
	// @Description サーバの稼働確認のためのエンドポイント
	// @Tags health
	// @Produce plain
	// @Success 200 {string} string "OK"
	// @Router / [get]
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// ゲーム開始エンドポイント
	// @Summary ゲーム開始
	// @Description 新しいゲームを開始する
	// @Tags game
	// @Accept json
	// @Produce json
	// @Success 201 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Router /games [post]
	r.POST("/games", gameServer.StartNewGame)

	// 判定エンドポイント
	// @Summary ゲームの推測
	// @Description ゲームに対して推測を送信し、判定結果を返す
	// @Tags game
	// @Accept json
	// @Produce json
	// @Param gameId path string true "ゲームID"
	// @Param guess body server.GuessRequest true "推測情報"
	// @Success 200 {object} server.GuessResponse
	// @Failure 400 {object} map[string]interface{}
	// @Failure 404 {object} map[string]interface{}
	// @Router /games/{gameId}/guess [post]
	r.POST("/games/:gameId/guess", gameServer.MakeGuess)

	log.Println("サーバーを開始します: http://localhost:8080")
	r.Run(":8080")
}
