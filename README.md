WsTestTool
===

# Usage

### Start with go: go run main.go -h
### Start with os version: ./exe/${OS_VERSION} -h

```bash
# ./exe/mac_amd64 -h

  -H string
    	Host.
  -d int
    	Time duration between each request. (default 10)
  -h	Usage.
  -n int
    	Test times. (default 1)
  -r int
    	Re-send message times. (default 1)
  -req string
    	Json string param
  -to int
    	Max waitting time. (default 10)
  -w	Watch each resposnes.
```

# Simple test

```bash=
# ./exe/mac_amd64 -H "ws://echo.websocket.org" -n 2 -r 5 -w

{"command":"ping"}
{"command":"ping"}
{"command":"ping"}
{"command":"ping"}
{"command":"ping"}
{"command":"ping"}
{"command":"ping"}
{"command":"ping"}
{"command":"ping"}
{"command":"ping"}

============================Sent Done===================================
INFO[0001] [Execution delay between repeat: 10ns]
INFO[0001] [Number of concurrent connections: 2]
INFO[0001] [Number of successful requests: 10]
INFO[0001] [Time of connections establishment: 954.477849ms]
INFO[0001] [Total execution time: 1.190156998s]
INFO[0001] [Average response time: 119.015699ms]
WARN[0001] [Number of failed requests: 0]
=======================================================================
```