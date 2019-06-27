package kernel

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

//Conn connection&lock
type Conn struct {
	Conn *websocket.Conn
	mu   *sync.Mutex
}

//SendMsg Send msg with lock controller
func (ws *Conn) sendMsg(msg []byte) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	err := ws.Conn.WriteMessage(1, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

