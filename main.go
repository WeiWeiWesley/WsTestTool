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

//統計&核心
var (
	successSignal     = make(chan bool)
	sentSignal        = make(chan bool)
	closeSignal       = make(chan bool)
	success           int
	fail              int
	sentCount         int
	startTime         time.Time
	execTime          time.Duration
	avgTime           time.Duration
	maxTime           time.Duration
	sunTime           time.Duration
	endTimes          int
	connEstablishTime time.Duration
)

//參數
var (
	help           bool //使用方法
	threads        int  //併發數量
	delay          time.Duration
	d              int         //每筆發送延遲
	to             int         //最大等待秒數
	timeout        *time.Timer //最大等待時間
	host           string      //目標網域
	path           string      // 目標URL
	watch          bool        //觀察每筆回傳
	request        string      //json string param
	reqSend        []byte
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
	flag.IntVar(&threads, "n", 1, "Number of connections.")
	flag.IntVar(&d, "d", 10, "Time duration(nanosecond) between each request.")
	flag.IntVar(&to, "to", 10, "Max waitting time(second).")
	flag.IntVar(&repeat, "r", 1, "Re-send message times.")
	flag.StringVar(&request, "req", "", "Json string param")
	flag.StringVar(&timing, "timing", "", "Start at particular time ex. 2019-05-08 15:04:00")
	flag.Parse()

	delay = time.Duration(d)
	timeout = time.NewTimer(time.Duration(to) * time.Second)

	endTimes = threads
	if repeat > 1 {
		endTimes = threads * repeat
	}

	go func() {
		for {
			select {
			case <-sentSignal:
				sentCount++
			case <-successSignal:
				success++
				if success == endTimes {
					closeSignal <- true
				}
			}
		}
	}()

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

		log.Print("warn", "Expected start at: "+timing)
		waitTime := timeClock.Sub(now)
		timeout.Reset(waitTime + time.Duration(to)*time.Second)
		time.Sleep(waitTime)
	}

	//Sender
	run()

	//Watting result
	wait()
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
	startTime = time.Now()

	var (
		connPool = make(map[int]Conn) //連線池
	)

	//建立所有連線
	for i := 0; i < threads; i++ {
		conn, err := ws.Connect(host)
		if err != nil {
			log.Print("error", "Connection establish fail. Please check host")
			return
		}

		connPool[i] = Conn{
			Conn: conn,
			mu:   &sync.Mutex{},
		}
	}
	connEstablishTime = time.Since(startTime)

	//併發
	for i := 0; i < threads; i++ {
		//重複發送
		go func(ws Conn) {
			for i := 0; i < repeat; i++ {
				err := ws.sendMsg(reqSend)
				if err != nil {
					fmt.Println("Send error:", err)
					return
				}

				sentSignal <- true
				time.Sleep(delay)
			}
		}(connPool[i])

		//個別監聽
		go func(ws Conn) {
			for {
				_, res, err := ws.Conn.ReadMessage()
				if err != nil {
					fmt.Println("Read error:", err)
					return
				}

				successSignal <- true
				if watch {
					fmt.Println(string(res))
				}
			}
		}(connPool[i])
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

	if request != "" {
		res, err := json.Marshal(request)
		if err != nil {
			return errors.New("Please check request format")
		}

		reqSend = res
	} else {
		msg := make(map[string]interface{})
		msg["command"] = "ping"
		b, err := json.Marshal(msg)
		if err != nil {
			return errors.New("Please check request format")
		}
		reqSend = b
	}

	return nil
}

func result(timeout bool) {
	msg := "Sent Done"
	if timeout {
		msg = "Timeout"
	}

	execTime = time.Since(startTime)
	avgTime = execTime / time.Duration(endTimes)

	fmt.Println()
	fmt.Println("============================" + msg + "===================================")
	log.Print("info", "Execution delay between repeat:", delay)
	log.Print("info", "Number of concurrent connections:", threads)
	log.Print("info", "Number of successful requests:", success)
	log.Print("info", "Time of connections establishment:", connEstablishTime)
	log.Print("info", "Total execution time:", execTime)
	log.Print("info", "Average response time:", avgTime)
	log.Print("warn", "Number of failed requests:", endTimes-success)
	fmt.Println("=======================================================================")
	time.Sleep(time.Second)
}
