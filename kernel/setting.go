package kernel

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"WsTestTool/log"
	"encoding/json"
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

}

//Help show usage
func Help() bool {
	return help
}

//Timer set timer
func Timer() {
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

	return
}

//Usage Show usage
func Usage() {
	flag.PrintDefaults()
}

//CheckParam Check param format
func CheckParam() (paramErr error) {
	if len(host) < 1 {
		paramErr = errors.New("Please use -H add host")
		log.Print("error", paramErr.Error())
	}

	if request != "" {
		res, err := json.Marshal(request)
		if err != nil {
			paramErr = errors.New("Please check request format")
			log.Print("error", paramErr.Error())
		} else {
			reqSend = res
		}

	} else {
		msg := make(map[string]interface{})
		msg["command"] = "ping"
		b, err := json.Marshal(msg)
		if err != nil {
			paramErr = errors.New("Please check request format")
			log.Print("error", paramErr.Error())
		} else {
			reqSend = b
		}
	}

	return
}
