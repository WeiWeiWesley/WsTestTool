package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"

	"WsTestTool/log"
	"WsTestTool/ws"

	"github.com/gorilla/websocket"
)

//Receive 接收範例
type Receive struct {
	Cmd   string `json:"commond"`
	Code  int    `json:"code"`
	Error string `json:"error"`
}

//統計&核心
var (
	closeSignal = make(chan bool)
	success     int
	fail        int
	sentCount   int
	startTime   time.Time
	execTime    time.Duration
	avgTime     time.Duration
	maxTime     time.Duration
	sunTime     time.Duration
	endTimes    int
)

//參數
var (
	help           bool //使用方法
	times          int  //測試次數
	delay          time.Duration
	d              int         //每筆發送延遲
	to             int         //最大等待秒數
	timeout        *time.Timer //最大等待時間
	host           string      //目標網域
	path           string      // 目標URL
	watch          bool        //觀察每筆回傳
	request        string      //json string param
	repeat         int
	repeatDuration time.Duration
	timing         string
)

func init() {
	//Require
	flag.StringVar(&host, "H", "", "Host.")
	//Options
	flag.BoolVar(&help, "h", false, "Usage.")
	flag.BoolVar(&watch, "w", false, "Watch each resposnes.")
	flag.IntVar(&times, "n", 1, "Test times.")
	flag.IntVar(&d, "d", 10, "Time duration between each request.")
	flag.IntVar(&to, "to", 10, "Max waitting time.")
	flag.IntVar(&repeat, "r", 1, "Re-send message times.")
	flag.StringVar(&request, "req", "", "Json string param")
	flag.StringVar(&timing, "timing", "", "Start at particular time ex. 2019-05-08 15:04:00")
	flag.Parse()

	delay = time.Duration(d)
	timeout = time.NewTimer(time.Duration(to) * time.Second)

	endTimes = times
	if repeat > 1 {
		endTimes = times * repeat
	}

}

func main() {
	if help {
		flag.Usage()
		return
	}

	//參數檢驗
	if err := checkParam(); err != nil {
		log.Print("error", err.Error())
		fmt.Println()
		flag.Usage()
		return
	}

	//定時器
	if len(timing) > 0 {
		timeClock, err := time.Parse("2006-01-02 15:04:05", timing)
		if err != nil {
			fmt.Println(err)
			return
		}

		localTime, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println(err)
			return
		}
		timeClock = timeClock.In(localTime).Add(-8 * time.Hour)

		now := time.Now()
		if now.After(timeClock) {
			log.Print("error", "Time "+timing+" should after current time.")
			fmt.Println()
			flag.Usage()
			return
		}

		log.Print("warn", "預計於: "+timing+" 開始執行...")
		waitTime := timeClock.Sub(now)
		timeout.Reset(waitTime + time.Duration(to)*time.Second)
		time.Sleep(waitTime)
	}

	//Sender
	run()

}

//Conn Conn
type Conn struct {
	Conn *websocket.Conn
	mu   *sync.Mutex
}

//SendMsg SendMsg
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

//Sender
func run() {
	ws.Init()
	startTime = time.Now()

	conn, err := ws.Connect(host)
	if err != nil {
		fmt.Println(err)
	}

	wsConn := Conn{
		Conn: conn,
		mu:   &sync.Mutex{},
	}


	msg := make(map[string]interface{})
	msg["command"] = "ping"

	b, _ := json.Marshal(msg)

	wsConn.sendMsg(b)

	for {
		_, res, err := wsConn.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(res))
	}
}

//Receiver
func wait() {
	select {
	case <-timeout.C:
		result(true)
	case <-closeSignal:
		result(false)
	}
}

func usage() {
	flag.PrintDefaults()
}

func checkParam() error {
	if len(host) < 1 {
		return errors.New("Please use -H add host")
	}

	return nil
}

func result(timeout bool) {
	msg := "發送完畢"
	if timeout {
		msg = "等待逾時"
	}

	execTime = time.Since(startTime)
	avgTime = sunTime / time.Duration(endTimes)

	fmt.Println()
	fmt.Println("============================" + msg + "===================================")
	log.Print("info", "執行延遲:", delay, "; 1 nanosecond = 0.0000000001 seconds")
	log.Print("info", "併發連線數量:", times)
	log.Print("info", "成功請求數量:", success)
	log.Print("info", "總執行時間:", execTime)
	log.Print("info", "平均回應時間:", avgTime)
	log.Print("info", "最大回應時間:", maxTime)
	log.Print("warn", "失敗請求數量:", endTimes-success)
	fmt.Println("=======================================================================")
	time.Sleep(time.Second)
}
