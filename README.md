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
# ./exe/mac_amd64 -H "ws://echo.websocket.org" -w -r 3 -n 3

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
INFO[0001] [Number of concurrent connections: 3]
INFO[0001] [Number of successful requests: 9]
INFO[0001] [Total execution time: 1.61317515s]
INFO[0001] [Average response time: 179.241683ms]
WARN[0001] [Number of failed requests: 0]
=======================================================================
```