/**
 * Copyright 2025 dqx0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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

	// @Summary ヘルスチェック
	// @Description サーバの稼働確認のためのエンドポイント
	// @Tags health
	// @Produce plain
	// @Success 200 {string} string "OK"
	// @Router / [get]
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/games", gameServer.StartNewGame)

	r.POST("/games/:gameId/guess", gameServer.MakeGuess)

	log.Println("サーバーを開始します: http://localhost:8080")
	r.Run(":8080")
}
