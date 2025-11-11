package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type DeckResponse struct {
	Success   bool   `json:"success"`
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type NewCardResponse struct {
	Success bool   `json:"success"`
	DeckID  string `json:"deck_id"`
	Cards   []struct {
		Code  string `json:"code"`
		Image string `json:"image"`
		Value string `json:"value"`
		Suit  string `json:"suit"`
	} `json:"cards"`

	Remaining int `json:"remaining"`
}

func generateDeck() string {
	resp, err := http.Get("https://deckofcardsapi.com/api/deck/new/shuffle/?deck_count=1")
	if err != nil {
		return "Error fetching deck"
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response"
	}

	var deckResponse DeckResponse
	err = json.Unmarshal(body, &deckResponse)
	if err != nil {
		return "Error parsing JSON"
	}

	return deckResponse.DeckID
}

func generateCardValue(deckID string) string {
	url := fmt.Sprintf("https://deckofcardsapi.com/api/deck/%s/draw/?count=1", deckID)

	resp, err := http.Get(url)
	if err != nil {
		return "Error fetching card"
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response"
	}

	var cardResponse NewCardResponse
	err = json.Unmarshal(body, &cardResponse)
	if err != nil {
		return "Error parsing JSON"
	}

	if len(cardResponse.Cards) > 0 {
		return cardResponse.Cards[0].Value
	}

	panic("no cards returned")
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	deckID := generateDeck()

	fmt.Println("Твой предикт, какой по счету выпадет дама: ")
	input, _ := reader.ReadString('\n')

	predict, err := strconv.Atoi(input[:len(input)-1])
	if err != nil {
		fmt.Println("Ты конечно офигел, надо ввести число.")
		return
	}

	if predict <= 0 || predict > 52 {
		fmt.Println("Матчасть учи, число должно быть от 1 до 52.")
		return
	}

	for i := 1; ; i++ {
		cardValue := generateCardValue(deckID)
		fmt.Printf("Выпала карта: %s\n", cardValue)
		if cardValue == "QUEEN" {
			fmt.Println("Выпала дама!")
			if predict == i {
				fmt.Println("Красава, ты угадал, держи бипки!")
			} else {
				fmt.Printf("Ну ты и лох, дама то выпала под номером %d, а не %d\n", i, predict)
			}

			break
		}

		if i == predict {
			fmt.Printf("Ну ты и нуб, я сгенерировал уже %d карт, а дамы все нет!", predict)
			return
		}
	}
}
