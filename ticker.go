package main

import (
	"encoding/json"
	"github.com/toorop/go-pusher"
	"log"
	"time"
)

const BitstampAppKey = "de504dc5763aeef9ff52"

type BitstampEventStub struct {
	Price float64 `json:"price"`
}

type Ticker struct{}

func (t *Ticker) Start(priceChannel chan float64, disconnectChannel chan bool) {
	for {
		t.ConnectPusher(priceChannel)
		disconnectChannel <- true
		time.Sleep(5 * time.Second)
	}
}

// Function returns on error
func (t *Ticker) ConnectPusher(priceChannel chan float64) {
	// Connect
	pusherClient, err := pusher.NewClient(BitstampAppKey)
	if err != nil {
		log.Println("Connection error: ", err)
		return
	}
	defer pusherClient.Close()
	log.Println("Connected to Pusher.")

	// Subscribe
	err = pusherClient.Subscribe("live_trades")
	if err != nil {
		log.Println("Subscription error: ", err)
		return
	}

	// Bind channels
	tradeChannel, err := pusherClient.Bind("trade")
	if err != nil {
		log.Println("Bind error: ", err)
		return
	}
	errChannel, err := pusherClient.Bind(pusher.ErrEvent)
	if err != nil {
		log.Println("Bind error: ", err)
		return
	}

	for {
		select {
		case tradeEvent := <-tradeChannel:
			var bitstampEvent BitstampEventStub
			err = json.Unmarshal([]byte(tradeEvent.Data), &bitstampEvent)
			if err != nil {
				log.Println("JSON error: ", err)
				return
			}
			priceChannel <- bitstampEvent.Price
		case errEvent := <-errChannel:
			log.Println("ErrEvent: " + errEvent.Data)
			return
		}
	}
}
