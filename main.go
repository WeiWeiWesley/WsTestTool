package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"

	"WsTestTool/log"
	"WsTestTool/ws"
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

	//Sender
	go run()

	//Waitting response
	wait()
}

//Sender
func run() {
	ws.Init()
	go monitor()

	startTime = time.Now()
	for i := 0; i < times; i++ {
		time.Sleep(delay * time.Nanosecond)
		if err := ws.Connect(strconv.Itoa(i), host, path, repeat); err != nil {
			fail++
			continue
		}

		msg := make(map[string]interface{})
		if request != "" {
			err := json.Unmarshal([]byte(request), &msg)
			if err != nil {
				log.Print("error", err.Error())
				closeSignal <- true
			}
		} else {
			msg["command"] = "ping"
		}

		msg["key"] = strconv.Itoa(i)

		ws.Send(msg)
	}

	go func() {
		if fail == endTimes {
			closeSignal <- true
		}
	}()
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

func monitor() {
	for {
		select {
		case data := <-ws.ReceiveChan:
			sentCount++
			response := string(data.Message)

			if watch {
				log.Print("info", "Sent count:", sentCount, "Response:", response, "Time:", data.TimeSpent.String())
			}

			if len(response) > 0 {
				success++
			} else {
				fail++
			}

			if maxTime < data.TimeSpent {
				maxTime = data.TimeSpent
			}

			sunTime += data.TimeSpent

			if sentCount == endTimes {
				closeSignal <- true
			}
		}
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
