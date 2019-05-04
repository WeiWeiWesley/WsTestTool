package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type ConnInfo struct {
	Ws       *websocket.Conn
	SendTime time.Time
}

//連線
var (
	Conn map[string]ConnInfo
	sendChan chan map[string]interface{}
)

//Init 設定初始化
func Init() {
	//Channel & map init
	Conn = make(map[string]ConnInfo)

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
					Conn[key] = ConnInfo{
						Ws: conn.Ws,
						SendTime: time.Now(),
					}
					err := conn.Ws.WriteMessage(websocket.TextMessage, msg)
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
