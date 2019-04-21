package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

//連線
var (
	Conn   map[string]*websocket.Conn
	sendChan chan map[string]interface{}
)

//Init 設定初始化
func Init() {
	//Channel & map init
	Conn = make(map[string]*websocket.Conn)
	sendChan = make(chan map[string]interface{})
	ReceiveChan = make(chan Receive)
}

//Connect Add Connection setting && listening
func Connect(key, host, path string) error {
	//連線 逾時 3s
	websocket.DefaultDialer.HandshakeTimeout = 3 * time.Second
	if err := ConnServer(key, host, path); err != nil {
		return err
	}

	keepWS(key)

	return nil
}

//Start listening
func keepWS(key string) {
	//Send message
	go func() {
		for {
			select {
			case data := <-sendChan:
				key, ok := data["key"].(string)
				if !ok {
					fmt.Println(data)
					break
				}

				delete(data, "key")
				msg, _ := json.Marshal(data)
				
				//檢查map連線存在
				if conn, ok := Conn[key]; ok {
					err := conn.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						fmt.Println(err)
					}

				} else {
					fmt.Println("Connection not exists")
				}
			}
		}
	}()

	//Receive message
	go receive(key) //websocket receiver
}
