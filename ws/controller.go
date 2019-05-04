package ws

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
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

//ConnServer 與 Server建立連線
func ConnServer(key, host, path string) error {
	jpURL := url.URL{Scheme: "ws", Host: host, Path: path}
	conn, _, err := websocket.DefaultDialer.Dial(jpURL.String(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	Conn[key] = ConnInfo{Ws: conn}

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
