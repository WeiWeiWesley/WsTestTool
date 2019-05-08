package ws

import (
	"fmt"
	"time"
)

//Receive Receive data
type Receive struct {
	Error     error
	Message   []byte
	Key       string
	TimeSpent time.Duration
}

//ReceiveChan Jackpot channel
var ReceiveChan chan Receive

//Send Api發送
func Send(data map[string]interface{}) error {
	//發送目標
	go func() {
		sendChan <- data
	}()

	return nil
}

//Websocket receiver
func receive(key string) {
	//等待回傳
	for {
		if _, ok := Conn[key]; ok {
			_, message, err := Conn[key].Ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			// fmt.Println(key, string(message)) //DEBUG
			ReceiveChan <- Receive{
				Error:   err,
				Message: message,
				Key:     key,
				TimeSpent: time.Since(Conn[key].SendTime),
			}
		}
	}
}
