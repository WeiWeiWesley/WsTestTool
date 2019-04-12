package main

import (
	"fmt"
	"os"
	"strconv"
	"testJpServer/ws"
	"time"
)

const (
	times = 100
)

func main() {
	env := os.Getenv("ENV")
	cf := "config/env/" + env + "_config.toml"

	go getJpChan()

	for i := 0; i < times; i++ {
		time.Sleep(10)
		ws.Init(cf, strconv.Itoa(i))
	}

	time.Sleep(2 * time.Second)
}

func getJpChan() {
	for {
		select {
		case data := <-ws.JPChan:
			fmt.Println(string(data.Message))
		}
	}
}
