package main

import "fmt"

func main() {
	// Make channels
	disconnectChannel := make(chan bool)
	priceChannel := make(chan float64)
	ledChannel := make(chan LedArray)

	// Init board
	board := Board{}
	go board.Start(ledChannel)

	// Init price ticker
	ticker := Ticker{}
	go ticker.Start(priceChannel, disconnectChannel)

	// Make "loading" screen
	ledArrayIdle := makeLoadingScreen()
	ledChannel <- *ledArrayIdle

	text := NewText()
	var lastPrice string
	for {
		select {
		case price := <-priceChannel:
			newPrice := fmt.Sprintf("%.0f", price)
			if newPrice != lastPrice {
				ledArray := text.MakeText(newPrice)
				ledChannel <- *ledArray
				lastPrice = newPrice
			}

		case <-disconnectChannel:
			ledChannel <- *ledArrayIdle
		}
	}
}

func makeLoadingScreen() *LedArray {
	ledArray := LedArray{}
	for x := 8; x <= 12; x += 2 {
		ledArray[x][6] = true
	}
	return &ledArray
}
