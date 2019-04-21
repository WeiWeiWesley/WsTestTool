package main

import (
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
)

//參數
var (
	help    bool //使用方法
	times   int  //測試次數
	delay   time.Duration
	d       int    //每筆發送延遲
	timeout int    //最大等待時間
	host    string //目標網域
	path    string // 目標URL
	watch   bool   //觀察每筆回傳
)

func init() {
	//Require
	flag.StringVar(&host, "H", "", "Host.")
	flag.StringVar(&path, "P", "", "URL path.")
	//Options
	flag.BoolVar(&help, "h", false, "Usage.")
	flag.BoolVar(&watch, "w", false, "Watch each resposnes.")
	flag.IntVar(&times, "n", 1, "Test times.")
	flag.IntVar(&d, "d", 10, "Time duration between each request.")
	flag.IntVar(&timeout, "timeout", 10, "Test times.")
	flag.Parse()

	delay = time.Duration(d)
}

func main() {
	//參數檢驗
	if err := checkParam(); err != nil {
		log.Print("error", err.Error())
		fmt.Println()
		flag.Usage()
		return
	}

	//Sender
	run()

	//Receiver
	result()
}

//Sender
func run() {
	ws.Init()
	go countResult()

	for i := 0; i < times; i++ {
		time.Sleep(delay * time.Nanosecond)
		if err := ws.Connect(strconv.Itoa(i), host, path); err != nil {
			fail++
			continue
		}

		msg := make(map[string]interface{})
		msg["command"] = "ping"
		msg["key"] = strconv.Itoa(i)

		ws.Send(msg)
	}

	go func() {
		if fail == times {
			closeSignal <- true
		}
	}()
}

//Receiver
func result() {
	select {
	case <-closeSignal:
		fmt.Println()
		fmt.Println("=======================================================================")
		log.Print("info", "執行數量:", times)
		log.Print("info", "執行延遲:", delay, "; 1 nanosecond = 0.0000000001 seconds")
		log.Print("info", "成功數量:", success)
		log.Print("info", "失敗數量:", fail)
		log.Print("info", "發送完畢")
		fmt.Println("=======================================================================")
		time.Sleep(time.Second)
	}
}

func countResult() {
	for {
		select {
		case data := <-ws.ReceiveChan:
			sentCount++
			response := string(data.Message)

			if watch {
				log.Print("info", "Sent count:", sentCount, "Response:", response)
			}

			if len(response) > 0 {
				success++
			} else {
				fail++
			}

			if sentCount == times {
				closeSignal <- true
			}
		}
	}
}

func usage() {
	flag.PrintDefaults()
}

func checkParam() error {
	if help {
		flag.Usage()
		return errors.New("")
	}

	if len(host) < 1 {
		return errors.New("Please use -H add host")
	}

	if len(path) < 1 {
		return errors.New("Please use -P add url path")
	}

	return nil
}
