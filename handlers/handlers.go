package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/malikov0216/websocket/models"
)

var (
	baseURL         = "wss://stream.binance.com:9443/ws"
	pingPongTimeout = time.Second * 60
)

// DepthWebSocket function that handles message from client
func DepthWebSocket(message []byte) {
	depth := new(models.Depth)
	if err := json.Unmarshal(message, &depth); err != nil {
		log.Println(err)
	}
	if depth.Event != "" && len(depth.Asks) != 0 && len(depth.Bids) != 0 {
		_, err := fmt.Printf("{bid: {price:%s, amount:%s}, ask: {price:%s, amount:%s}}\n", depth.Bids[0][0], depth.Bids[0][1], depth.Asks[0][0], depth.Asks[0][1])
		if err != nil {
			log.Println(err)
		}
	}
}

// ErrorWebSocket defines error
func ErrorWebSocket(err error) {
	log.Println(err)
}

// SubscribeWebSocket subscribe to stream
func SubscribeWebSocket(id int, symbol string, handler models.DepthWebSocket, errHandler models.ErrorWebSocket) (closeChanel, unSubcribe chan struct{}, err error) {
	connection, _, err := websocket.DefaultDialer.Dial(baseURL, nil)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("ID:", id)
	param := fmt.Sprintf("%s@depth", symbol)

	subscribeMap := models.LiveStream{"SUBSCRIBE", []string{param}, id}
	unSubcribeMap := models.LiveStream{"UNSUBSCRIBE", []string{param}, id}

	err = connection.WriteJSON(subscribeMap)
	if err != nil {
		return nil, nil, err
	}

	closeChanel = make(chan struct{})
	unSubcribe = make(chan struct{})

	go chanelControl(connection, unSubcribeMap, handler, errHandler, closeChanel, unSubcribe)
	return
}

func chanelControl(connection *websocket.Conn, unSubcribeMap models.LiveStream, handler models.DepthWebSocket, errHandler models.ErrorWebSocket, closeChanel, unSubcribe chan struct{}) {
	defer func() {
		connectionError := connection.Close()
		if connectionError != nil {
			errHandler(connectionError)
		}
	}()
	defer close(closeChanel)
	sendPongFrame(connection, pingPongTimeout)

	for {
		select {
		case <-unSubcribe:
			err := connection.WriteJSON(unSubcribeMap)
			if err != nil {
				errHandler(err)
				return
			}
			return
		default:
			_, message, err := connection.ReadMessage()
			if err != nil {
				errHandler(err)
				return
			}
			handler(message)
		}
	}
}

func sendPongFrame(connection *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	connection.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := connection.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Now().Sub(lastResponse) > timeout {
				connection.Close()
				return
			}
		}
	}()
}
