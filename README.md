WsTestTool
===

# Usage   

### Start with go: go run main.go -h
### Start with os version: ./exe/${OS_VERSION} -h

```
# ./exe/mac_amd64 -h

  -H string
    	Host.
  -d int
    	Time duration(nanosecond) between each request (default 10)
  -h	Usage
  -n int
    	Number of connections (default 1)
  -r int
    	Re-send message times (default 1)
  -req string
    	Json string param
  -resEq string
    	Check response msg equivalent to something
  -resHas string
    	Check response msg contain sub string
  -timing string
    	Start at particular time ex. 2019-05-08 15:04:00
  -to int
    	Max waitting time(second) (default 10)
  -w	Watch each resposnes
```

# Simple Start

```bash
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

# Response Verification

## Contains string

#### Success

```bash
# ./exe/mac_amd64 -H "ws://echo.websocket.org" -w -resHas 'ping'

{"command":"ping"}

============================Sent Done===================================
INFO[0000] [Execution delay between repeat: 10ns]
INFO[0000] [Number of concurrent connections: 1]
INFO[0000] [Number of successful requests: 1]
INFO[0000] [Time of connections establishment: 431.996841ms]
INFO[0000] [Total execution time: 647.519532ms]
INFO[0000] [Average response time: 647.519532ms]
WARN[0000] [Number of failed requests: 0]
=======================================================================
```

#### Fail
```bash
# ./exe/mac_amd64 -H "ws://echo.websocket.org" -w -resHas 'wwww'

String: "wwww" not in Response:"{"command":"ping"}"

============================Timeout===================================
INFO[0010] [Execution delay between repeat: 10ns]
INFO[0010] [Number of concurrent connections: 1]
INFO[0010] [Number of successful requests: 0]
INFO[0010] [Time of connections establishment: 460.937178ms]
INFO[0010] [Total execution time: 9.999936605s]
INFO[0010] [Average response time: 9.999936605s]
WARN[0010] [Number of failed requests: 1]
=======================================================================
```

## Same string

#### Success

```bash
# ./exe/mac_amd64 -H "ws://echo.websocket.org" -w -resEq '{"command":"ping"}'

{"command":"ping"}

============================Sent Done===================================
INFO[0000] [Execution delay between repeat: 10ns]
INFO[0000] [Number of concurrent connections: 1]
INFO[0000] [Number of successful requests: 1]
INFO[0000] [Time of connections establishment: 531.191969ms]
INFO[0000] [Total execution time: 757.869264ms]
INFO[0000] [Average response time: 757.869264ms]
WARN[0000] [Number of failed requests: 0]
=======================================================================
```

#### Fail

```bash
# ./exe/mac_amd64 -H "ws://echo.websocket.org" -w -resEq 'ping'

resEq: "ping" != Response: "{"command":"ping"}"

============================Timeout===================================
INFO[0010] [Execution delay between repeat: 10ns]
INFO[0010] [Number of concurrent connections: 1]
INFO[0010] [Number of successful requests: 0]
INFO[0010] [Time of connections establishment: 590.022344ms]
INFO[0010] [Total execution time: 10.004899823s]
INFO[0010] [Average response time: 10.004899823s]
WARN[0010] [Number of failed requests: 1]
=======================================================================
```
