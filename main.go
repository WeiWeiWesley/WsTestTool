package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"WsTestTool/log"
	"WsTestTool/ws"
)

type Receive struct {
	Cmd   string `json:"commond"`
	Code  int    `json:"code"`
	Error string `json:"error"`
}

var (
	closeSignal = make(chan bool)
	success     int
	fail        int
	sentCount   int
)

var (
	config  string
	help    bool
	times   int
	delay   time.Duration
	d       int
	timeout int
)

func init() {
	flag.BoolVar(&help, "h", false, "Usage.")
	flag.StringVar(&config, "c", "local", "Test config file name.")
	flag.IntVar(&times, "n", 1, "Test times.")
	flag.IntVar(&d, "d", 10, "Time duration between each request.")
	flag.IntVar(&timeout, "timeout", 10, "Test times.")
}

func main() {
	flag.Parse()
	{
		delay = time.Duration(d)
	}

	if help {
		flag.Usage()
		return
	}

	run()

	result()
}

func run() {
	ws.Init("config/" + config + ".toml")
	go countResult()

	for i := 0; i < times; i++ {
		time.Sleep(delay * time.Nanosecond)
		ws.Connect(strconv.Itoa(i))
	}
}

func result() {
	select {
	case <-closeSignal:
		fmt.Println()
		log.Print("info", "執行數量:", times)
		log.Print("info", "執行延遲:", delay, "; 1 nanosecond = 0.0000000001 seconds")
		log.Print("info", "成功數量:", success)
		log.Print("info", "失敗數量:", fail)
		log.Print("info", "發送完畢")
		time.Sleep(time.Second)
	}
}

func countResult() {
	for {
		select {
		case data := <-ws.JPChan:
			sentCount++
			fmt.Println("count", success, string(data.Message))

			var result Receive
			json.Unmarshal(data.Message, &result)
			if result.Code != 200 {
				fail++
			} else {
				success++
			}

			if sentCount == times {
				closeSignal <- true
			}
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `
		Usage: Websocket Test Tool
		Options:
	`)

	flag.PrintDefaults()
}
