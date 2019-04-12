package ws

import (
	"encoding/json"
	"fmt"
	"net/url"
	"starfruit/kernel/common/log"

	"github.com/gorilla/websocket"
)

type Receive struct {
	Error   error
	Message []byte
}

//JPChan Jackpot channel
var JPChan chan Receive
var jpConnToken string

//SendToJackpot Api發送
func SendToJackpot(data map[string]interface{}) error {
	//發送目標
	go func() {
		data["auth"] = "authorization"
		msg, _ := json.Marshal(data)
		jpSend <- msg
	}()

	return nil
}

//ConnJpServer 與 JP Server建立連線
func ConnJpServer(reTry bool, token, key string) {
	jpURL := url.URL{Scheme: "ws", Host: config.API.JP.Host + ":" + config.API.JP.Port, Path: "ws/keep"}
	jpC, _, err := websocket.DefaultDialer.Dial(jpURL.String(), nil)
	if err != nil {
		log.Print("error", fmt.Sprintf("[Jackpot] Error dial: %+v url: %v", err, jpURL.String()))
		return
	}

	Conn["jackpot_server"+key] = jpC

	var regist = make(map[string]interface{})
	if reTry {
		go receiveJackpot(key)
		regist["command"] = "server_reconnection"
		regist["data"] = `{"token":"` + token + `"}`
	} else {
		regist["command"] = "server_regist"
		regist["data"] = `{"ip":"127.0.0.1"}`
	}

	SendToJackpot(regist)
}

//Jackpot websocket receiver
func receiveJackpot(key string) {
	//等待回傳
	for {
		if _, ok := Conn["jackpot_server"+key]; ok {
			_, message, err := Conn["jackpot_server"+key].ReadMessage()
			if err != nil {
				log.Print("error", "Jackpot server connection fail: "+err.Error())
				// reConnection()
				return
			}

			fmt.Printf("response[%s]: %+v \n", key, string(message))

			JPChan <- Receive{
				Error:   err,
				Message: message,
			}
		}
	}
}

//ReConnection re-connection to JP server
func ReConnection(token string, key string) {
	log.Print("warn", "Try re-connection to JP server use token: "+token)
	ConnJpServer(true, token, key)
}

//SetToken SetToken
func SetToken(token string) {
	jpConnToken = token
}

//GetToken GetToken
func GetToken() string {
	return jpConnToken
}
