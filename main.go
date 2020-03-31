package main

import (
	"log"
	"time"

	"github.com/malikov0216/websocket/handlers"
)

func main() {
	var errOut error
	closeChanel, unSubscribe, errOut := handlers.SubscribeWebSocket(1, "btcusdt", handlers.DepthWebSocket, handlers.ErrorWebSocket)
	closeChanel2, unSubscribe2, errOut := handlers.SubscribeWebSocket(2, "ethusdt", handlers.DepthWebSocket, handlers.ErrorWebSocket)
	closeChanel3, unSubscribe3, errOut := handlers.SubscribeWebSocket(3, "eosusdt", handlers.DepthWebSocket, handlers.ErrorWebSocket)
	if errOut != nil {
		log.Println(errOut)
		return
	}

	go startStream(unSubscribe, 5)
	go startStream(unSubscribe2, 5)
	go startStream(unSubscribe3, 5)

	<-closeChanel
	<-closeChanel2
	<-closeChanel3
}

func startStream(unSubscribe chan struct{}, td time.Duration) {
	if td != 0 {
		time.Sleep(td * time.Second)
	}
	unSubscribe <- struct{}{}
}
