package server

import (
	"crypto/rand"
	"net/http"
	"sync"

	"github.com/dqx0/hitandblow/internal/game"

	"github.com/gin-gonic/gin"
)

type GameServer struct {
	games map[string]*game.Game
	mu    sync.RWMutex
}

type GuessRequest struct {
	Number string `json:"number"`
}

type GuessResponse struct {
	Hit   int  `json:"hit"`
	Blow  int  `json:"blow"`
	Tries int  `json:"tries"`
	Clear bool `json:"clear"`
}

type StartNewGameResponse struct {
	GameID  string `json:"game_id"`
	Message string `json:"message"`
}

func NewGameServer() *GameServer {
	return &GameServer{
		games: make(map[string]*game.Game),
	}
}

// @Summary ゲーム開始
// @Description 新しいゲームを開始し、ゲームIDとメッセージを返す
// @Tags game
// @Accept json
// @Produce json
// @Success 201 {object} server.StartNewGameResponse
// @Router /games [post]
func (s *GameServer) StartNewGame(c *gin.Context) {
	s.mu.Lock()
	gameID := generateGameID()
	s.games[gameID] = game.NewGame()
	s.mu.Unlock()

	resp := StartNewGameResponse{
		GameID:  gameID,
		Message: "新しいゲームを開始しました",
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary ゲームの推測
// @Description ゲームに対して推測を送信し、判定結果を返す
// @Tags game
// @Accept json
// @Produce json
// @Param gameId path string true "ゲームID"
// @Param guess body server.GuessRequest true "推測情報"
// @Success 200 {object} server.GuessResponse "推測結果"
// @Failure 400 {object} object "無効なリクエスト"
// @Failure 404 {object} object "ゲームが見つかりません"
// @Router /games/{gameId}/guess [post]
func (s *GameServer) MakeGuess(c *gin.Context) {
	gameID := c.Param("gameId")

	s.mu.RLock()
	g, exists := s.games[gameID]
	s.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ゲームが見つかりません"})
		return
	}

	var req GuessRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト"})
		return
	}

	if !game.ValidateInput(req.Number) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な数字です"})
		return
	}

	result := g.Guess(req.Number)

	response := GuessResponse{
		Hit:   result.Hit,
		Blow:  result.Blow,
		Tries: g.GetTries(),
		Clear: result.Hit == 4,
	}

	c.JSON(http.StatusOK, response)
}

func generateGameID() string {
	return rand.Text()
}
