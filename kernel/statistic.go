package kernel

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"WsTestTool/log"
	"WsTestTool/ws"
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

func init() {
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

//Run Sender
func Run() {
	startTime = time.Now()

	var (
		connPool = make(map[int]Conn) //連線池
	)

	//建立所有連線
	for i := 0; i < threads; i++ {
		conn, err := ws.Connect(host)
		if err != nil {
			go func() {
				log.Print("error", "Connection establish fail. Please check host")
				closeSignal <- false
			}()
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
		go func(ws Conn, resEquivalent, resHas string) {
			for {
				_, res, err := ws.Conn.ReadMessage()
				if err != nil {
					fmt.Println("Read error:", err)
					return
				}

				resStr := string(res)
				if resEquivalent == "" && resHas == "" {
					successSignal <- true
					if watch {
						fmt.Println(resStr)
					}
					continue
				}

				if resEquivalent != "" {
					if resStr == resEquivalent {
						successSignal <- true
						if watch {
							fmt.Println(resStr)
						}
					} else {
						if watch {
							fmt.Printf("resEq: \"%s\" != Response: \"%s\" \n", resEquivalent, resStr)
						}
					}
					continue
				}

				if resHas != "" {
					if strings.Contains(resStr, resHas) {
						successSignal <- true
						if watch {
							fmt.Println(resStr)
						}
					} else {
						if watch {
							fmt.Printf("resHas: \"%s\" not in Response:\"%s\" \n", resHas, resStr)
						}
					}
					continue
				}
			}
		}(connPool[i], resEquivalent, resHas)
	}
}

//Wait Receiver
func Wait() {
	select {
	case <-timeout.C:
		result(true)
	case <-closeSignal:
		result(false)
	}
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
