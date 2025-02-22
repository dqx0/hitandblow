package game

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type Game struct {
	answer string
	tries  int
}

type Result struct {
	Hit  int
	Blow int
}

// @Summary ゲーム開始
// @Description 新しいゲームを開始する
// @Tags game
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /games [post]
func NewGame() *Game {
	return &Game{
		answer: generateAnswer(),
		tries:  0,
	}
}

func (g *Game) Guess(input string) Result {
	g.tries++
	hit := 0
	blow := 0

	usedAnswer := make([]bool, 4)
	usedInput := make([]bool, 4)

	for i := 0; i < 4; i++ {
		if input[i] == g.answer[i] {
			hit++
			usedAnswer[i] = true
			usedInput[i] = true
		}
	}

	for i := 0; i < 4; i++ {
		if usedInput[i] {
			continue
		}
		for j := 0; j < 4; j++ {
			if !usedAnswer[j] && input[i] == g.answer[j] {
				blow++
				usedAnswer[j] = true
				break
			}
		}
	}

	return Result{Hit: hit, Blow: blow}
}

func (g *Game) GetTries() int {
	return g.tries
}

func generateAnswer() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	result := ""
	for i := 0; i < 4; i++ {
		result += strconv.Itoa(numbers[i])
	}
	return result
}

func ValidateInput(input string) bool {
	if len(input) != 4 {
		return false
	}

	matched, _ := regexp.MatchString(`^\d{4}$`, input)
	if !matched {
		return false
	}

	used := make(map[rune]bool)
	for _, digit := range input {
		if used[digit] {
			return false
		}
		used[digit] = true
	}

	return true
}
