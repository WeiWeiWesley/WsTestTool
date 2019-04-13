package ws

import (
	"starfruit/kernel/common/log"
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
		log.Print("error", err.Error())
	}

	//Channel & map init
	Conn = make(map[string]*websocket.Conn)
	jpSend = make(chan []byte)
	JPChan = make(chan Receive)
}

//Connect Add Connection setting && listening
func Connect(key string) {
	//連線 逾時 3s
	websocket.DefaultDialer.HandshakeTimeout = 3 * time.Second
	//Jackpot
	ConnJpServer(false, "", key)
	keepWS(key)
}

//Start listening
func keepWS(key string) {
	//Send message
	go func(key string) {
		for {
			select {
			case msg := <-jpSend:
				//檢查map連線存在
				if conn, ok := Conn["jackpot_server"+key]; ok {
					err := conn.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						log.Print("error", "WriteMessage error: "+err.Error())

						//嘗試重新連線
						time.Sleep(5 * time.Second)
						ReConnection(jpConnToken, key)
					}

				} else {
					log.Print("error", "Conn[jackpot_server"+key+"] not exists")
				}
			}
		}
	}(key)

	//Receive message
	go receiveJackpot(key) //Jackpot websocket receiver
}
