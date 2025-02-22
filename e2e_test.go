package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dqx0/hitandblow/internal/server"
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
	gs := server.NewGameServer()
	r.POST("/games", gs.StartNewGame)
	r.POST("/games/:gameId/guess", gs.MakeGuess)
	return r
}

func TestE2EFlow(t *testing.T) {
	router := setupRouter()

	// ★ Step 1: ゲームを開始する
	reqStart, err := http.NewRequest(http.MethodPost, "/games", nil)
	if err != nil {
		t.Fatalf("新しいゲーム開始リクエストの作成に失敗: %v", err)
	}
	recStart := httptest.NewRecorder()
	router.ServeHTTP(recStart, reqStart)

	if recStart.Code != http.StatusCreated {
		t.Fatalf("ゲーム開始に失敗: HTTP %d", recStart.Code)
	}

	var startResp startResponse
	if err := json.Unmarshal(recStart.Body.Bytes(), &startResp); err != nil {
		t.Fatalf("ゲーム開始レスポンスのパースに失敗: %v", err)
	}
	if startResp.GameID == "" {
		t.Fatalf("game_id が空です")
	}

	// ★ Step 2: 推測リクエストを送信する
	// 有効な4桁かつ重複のない数字を使用
	guessPayload := `{"number": "1234"}`
	reqGuess, err := http.NewRequest(http.MethodPost, "/games/"+startResp.GameID+"/guess", bytes.NewBufferString(guessPayload))
	if err != nil {
		t.Fatalf("推測リクエストの作成に失敗: %v", err)
	}
	reqGuess.Header.Set("Content-Type", "application/json")
	recGuess := httptest.NewRecorder()
	router.ServeHTTP(recGuess, reqGuess)

	if recGuess.Code != http.StatusOK {
		t.Fatalf("推測リクエストが失敗: HTTP %d", recGuess.Code)
	}

	var guessResp guessResponse
	if err := json.Unmarshal(recGuess.Body.Bytes(), &guessResp); err != nil {
		t.Fatalf("推測レスポンスのパースに失敗: %v", err)
	}

	// 試行回数が1以上になっていることを確認
	if guessResp.Tries < 1 {
		t.Errorf("期待する試行回数は1以上ですが、実際は %d", guessResp.Tries)
	}
	// 答えはランダムなので、クリアしているかどうかは確認しない
	t.Logf("推測結果: %v", guessResp)
}
