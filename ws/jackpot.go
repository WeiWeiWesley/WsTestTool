package ws

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type Receive struct {
	Error   error
	Message []byte
	Key     string
}

//ReceiveChan Jackpot channel
var ReceiveChan chan Receive

//SendToJackpot Api發送
func SendToJackpot(data map[string]interface{}) error {
	//發送目標
	go func() {
		msg, _ := json.Marshal(data)
		jpSend <- msg
	}()

	return nil
}

//ConnJpServer 與 JP Server建立連線
func ConnJpServer(key, path string) error {
	var err error
	jpURL := url.URL{Scheme: "ws", Host: config.API.JP.Host + ":" + config.API.JP.Port, Path: path}
	Conn[key], _, err = websocket.DefaultDialer.Dial(jpURL.String(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

//Websocket receiver
func receive(key string) {
	//等待回傳
	for {
		if _, ok := Conn[key]; ok {
			_, message, err := Conn[key].ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			ReceiveChan <- Receive{
				Error:   err,
				Message: message,
				Key:     key,
			}
		}
	}
}
