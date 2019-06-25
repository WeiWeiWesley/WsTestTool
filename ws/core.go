package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//Connect Add Connection setting && listening
func Connect(host string) (*websocket.Conn, error) {
	//連線 逾時 3s
	websocket.DefaultDialer.HandshakeTimeout = 3 * time.Second
	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}
