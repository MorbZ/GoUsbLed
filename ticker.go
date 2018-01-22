package main

import (
	"encoding/json"
	"errors"
	"github.com/toorop/go-pusher"
	"log"
	"time"
)

const bitstampAppKey = "de504dc5763aeef9ff52"

type bitstampEventStub struct {
	Price float64 `json:"price"`
}

// Ticker informs about price changes via channels
type Ticker struct{}

// Start connection with Pusher
func (t *Ticker) Start(priceChannel chan float64, disconnectChannel chan bool) {
	for {
		t.connectPusher(priceChannel)
		disconnectChannel <- true
		time.Sleep(5 * time.Second)
	}
}

// Function returns on error
func (t *Ticker) connectPusher(priceChannel chan float64) {
	// Connect
	pusherClient, err := pusher.NewClient(bitstampAppKey)
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
			price, err := getPriceForEventData(tradeEvent.Data)
			if err != nil {
				log.Println("JSON error: ", err)
				return
			}
			priceChannel <- price
		case errEvent := <-errChannel:
			log.Println("ErrEvent: " + errEvent.Data)
			return
		}
	}
}

func getPriceForEventData(eventData string) (float64, error) {
	var bitstampEvent bitstampEventStub
	err := json.Unmarshal([]byte(eventData), &bitstampEvent)
	if err != nil {
		return 0, err
	}
	if bitstampEvent.Price == 0 {
		err = errors.New("Price is 0")
		return 0, err
	}
	return bitstampEvent.Price, nil
}
