package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dqx0/hitandblow/internal/game"
	"github.com/gin-gonic/gin"
)

type startResponse struct {
	GameID  string `json:"game_id"`
	Message string `json:"message"`
}

type guessResponse struct {
	Hit   int  `json:"hit"`
	Blow  int  `json:"blow"`
	Tries int  `json:"tries"`
	Clear bool `json:"clear"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	gs := NewGameServer()
	r.POST("/games", gs.StartNewGame)
	r.POST("/games/:gameId/guess", gs.MakeGuess)
	return r
}

func TestStartNewGame(t *testing.T) {
	router := setupRouter()

	req, err := http.NewRequest(http.MethodPost, "/games", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var resp startResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if resp.GameID == "" {
		t.Errorf("Expected non-empty game_id")
	}
	if resp.Message != "新しいゲームを開始しました" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestMakeGuessGameNotFound(t *testing.T) {
	router := setupRouter()

	// 存在しないゲームIDでリクエスト
	url := "/games/nonexistent/guess"
	body := bytes.NewBufferString(`{"number": "1234"}`)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for not found game, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestMakeGuessInvalidJSON(t *testing.T) {
	router := setupRouter()

	// まず、新しいゲームを作成して game_id を取得
	reqStart, _ := http.NewRequest(http.MethodPost, "/games", nil)
	recStart := httptest.NewRecorder()
	router.ServeHTTP(recStart, reqStart)
	var resp startResponse
	if err := json.Unmarshal(recStart.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal start response: %v", err)
	}

	// リクエストボディに不正な JSON を設定
	url := "/games/" + resp.GameID + "/guess"
	body := bytes.NewBufferString(`invalid json`)
	reqGuess, _ := http.NewRequest(http.MethodPost, url, body)
	reqGuess.Header.Set("Content-Type", "application/json")
	recGuess := httptest.NewRecorder()
	router.ServeHTTP(recGuess, reqGuess)

	if recGuess.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, recGuess.Code)
	}
}

func TestMakeGuessValidRequest(t *testing.T) {
	// ゲーム作成用にサーバーインスタンスを直接操作することで、既知の状態に設定する例
	gs := NewGameServer()
	// ゲーム作成後、内部のゲームの答えはランダムですが、
	// 正常なリクエストの場合はレスポンスとして hit/blow などが返されるので、
	// ここではレスポンス構造自体が返ることを確認します。

	// ゲーム作成
	gameID := "test-game-id"
	// ※ ゲームの答えを直接書き換えられないため、ここでは通常の生成結果を利用
	gs.games[gameID] = game.NewGame()

	// gin コンテキストの生成
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// パラメーター設定
	c.Params = []gin.Param{{Key: "gameId", Value: gameID}}
	// リクエストボディは有効な JSON とする
	c.Request, _ = http.NewRequest(http.MethodPost, "/games/"+gameID+"/guess", bytes.NewBufferString(`{"number": "1234"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	gs.MakeGuess(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var gr guessResponse
	if err := json.Unmarshal(w.Body.Bytes(), &gr); err != nil {
		t.Fatalf("Failed to unmarshal guess response: %v", err)
	}
	// 応答内容が数値型で返っていることを確認（具体的な値はランダムのため検証しにくい）
	if gr.Tries < 1 {
		t.Errorf("Expected tries >= 1, got %d", gr.Tries)
	}
}
