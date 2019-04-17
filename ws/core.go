package ws

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
)

//Config 連線設定
type Config struct {
	API struct {
		JP struct {
			IP   string `toml:"ip"`
			Host string `toml:"host"`
			Port string `toml:"port"`
			Auth string `toml:"auth"`
		} `toml:"jackpot"`
	} `toml:"api"`
}

//連線
var (
	Conn   map[string]*websocket.Conn
	config Config
	jpSend chan []byte
)

//Init 設定初始化
func Init(path string) {
	//Load config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		fmt.Println(err)
	}

	//Channel & map init
	Conn = make(map[string]*websocket.Conn)
	jpSend = make(chan []byte)
	ReceiveChan = make(chan Receive)
}

//Connect Add Connection setting && listening
func Connect(key, path string) error {
	//連線 逾時 3s
	websocket.DefaultDialer.HandshakeTimeout = 3 * time.Second
	//Jackpot
	if err := ConnJpServer(key, path); err != nil {
		return err
	}
	
	keepWS(key)

	return nil
}

//Start listening
func keepWS(key string) {
	//Send message
	go func(key string) {
		for {
			select {
			case msg := <-jpSend:
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
	}(key)

	//Receive message
	go receive(key) //websocket receiver
}
