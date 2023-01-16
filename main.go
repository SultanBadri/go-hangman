package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Game struct {
	word           string
	guessed        string
	incorrect      int
	maxIncorrect   int
	lettersGuessed string
}

var words []string

func (g *Game) init() {
	rand.Seed(time.Now().UnixNano())
	g.word = words[rand.Intn(len(words))]
	g.guessed = strings.Repeat("_", len(g.word))
	g.incorrect = 0
	g.maxIncorrect = 6
}

func (g *Game) isWon() bool {
	return g.word == g.guessed
}

func (g *Game) isLost() bool {
	return g.incorrect >= g.maxIncorrect
}

func (g *Game) guess(char rune) {
    if strings.ContainsRune(g.lettersGuessed, char) {
        fmt.Println("You already guessed that letter")
				fmt.Println()
        return
    }
    g.lettersGuessed += string(char)
    if strings.ContainsRune(g.word, char) {
        for i, r := range g.word {
            if char == r {
                g.guessed = g.guessed[:i] + string(char) + g.guessed[i + 1:]
            }
        }
    } else {
        g.incorrect++
    }
		fmt.Println()
}

func (g *Game) display() {
	fmt.Println("Word:", g.guessed)
	fmt.Println("Incorrect Guesses:", g.incorrect)
  fmt.Println("Letters Guessed:", strings.Split(g.lettersGuessed, ""))
}

func main() {
	// get words from API
	resp, err := http.Get("https://raw.githubusercontent.com/dwyl/english-words/master/words_dictionary.json")
	if err != nil {
		fmt.Println("Error fetching words:", err)
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// parse json data
	var hashmap map[string]interface{}
	err = json.Unmarshal(data, &hashmap)
	if err != nil {
		fmt.Println("Error parsing json:", err)
		return
	}

	// convert json data to words slice
	for word := range hashmap {
		words = append(words, word)
	}

	var game Game
	game.init()
	var guess string
	for !game.isWon() && !game.isLost() {
		game.display()
		fmt.Print("Enter a letter: ")
		fmt.Scan(&guess)
		game.guess([]rune(guess)[0])
	}
	if game.isWon() {
		fmt.Println("You won!")
		fmt.Println("The correct word is:", game.word)
	} else {
		fmt.Println("You lost!")
		fmt.Println("The word was:", game.word)
	}
}
