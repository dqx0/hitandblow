package game

import "testing"

func TestValidateInput(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1234", true},
		{"1123", false},  // 重複あり
		{"12", false},    // 桁数不足
		{"abcd", false},  // 数字でない
		{"12345", false}, // 桁数オーバー
	}

	for _, tc := range tests {
		result := ValidateInput(tc.input)
		if result != tc.expected {
			t.Errorf("ValidateInput(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestNewGameAnswer(t *testing.T) {
	g := NewGame()
	// 答えは 4 桁の数字かつ各桁が重複していないこと
	answer := g.answer
	if len(answer) != 4 {
		t.Fatalf("NewGame answer length = %d, expected 4", len(answer))
	}
	used := make(map[rune]bool)
	for _, digit := range answer {
		if used[digit] {
			t.Errorf("Duplicate digit %c found in answer %s", digit, answer)
		}
		used[digit] = true
	}
}

func TestGuessFullHit(t *testing.T) {
	// テスト用に正解を固定
	g := NewGame()
	g.answer = "1234"

	result := g.Guess("1234")
	if result.Hit != 4 || result.Blow != 0 {
		t.Errorf("Guess full hit: expected (4 hit, 0 blow), got (%d hit, %d blow)", result.Hit, result.Blow)
	}
	if g.GetTries() != 1 {
		t.Errorf("GetTries expected 1, got %d", g.GetTries())
	}
}

func TestGuessPartial(t *testing.T) {
	// 正解: 1234
	g := NewGame()
	g.answer = "1234"

	// 入力例: "1325"
	// 比較:
	// 位置0: '1' == '1' -> hit
	// 位置1: '3' is not equal '2' ですが、answer[2]=='3' -> blow
	// 位置2: '2' is not equal '3' ですが、answer[1]=='2' -> blow
	// 位置3: '5' は正解内に存在しない
	result := g.Guess("1325")
	if result.Hit != 1 || result.Blow != 2 {
		t.Errorf("Guess partial: expected (1 hit, 2 blow), got (%d hit, %d blow)", result.Hit, result.Blow)
	}
	if g.GetTries() != 1 {
		t.Errorf("GetTries expected 1, got %d", g.GetTries())
	}
}

func TestGuessNoMatch(t *testing.T) {
	// 正解: 1234
	g := NewGame()
	g.answer = "1234"

	result := g.Guess("5678")
	if result.Hit != 0 || result.Blow != 0 {
		t.Errorf("Guess no match: expected (0 hit, 0 blow), got (%d hit, %d blow)", result.Hit, result.Blow)
	}
	if g.GetTries() != 1 {
		t.Errorf("GetTries expected 1, got %d", g.GetTries())
	}
}
