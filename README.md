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
# ./exe/mac_amd64 -H echo.websocket.org -w -r 3 -n 3

INFO[0000] [Sent count: 1 Response: {"command":"ping"} Time: 203.157448ms]
INFO[0000] [Sent count: 2 Response: {"command":"ping"} Time: 203.17124ms]
INFO[0000] [Sent count: 3 Response: {"command":"ping"} Time: 203.504778ms]
INFO[0001] [Sent count: 4 Response: {"command":"ping"} Time: 307.303975ms]
INFO[0001] [Sent count: 5 Response: {"command":"ping"} Time: 307.319302ms]
INFO[0001] [Sent count: 6 Response: {"command":"ping"} Time: 307.63367ms]
INFO[0001] [Sent count: 7 Response: {"command":"ping"} Time: 312.078172ms]
INFO[0001] [Sent count: 8 Response: {"command":"ping"} Time: 312.09139ms]
INFO[0001] [Sent count: 9 Response: {"command":"ping"} Time: 312.241335ms]

============================發送完畢===================================
INFO[0001] [執行延遲: 10ns ; 1 nanosecond = 0.0000000001 seconds]
INFO[0001] [併發連線數量: 3]
INFO[0001] [成功請求數量: 9]
INFO[0001] [總執行時間: 1.965507712s]
INFO[0001] [平均回應時間: 274.277923ms]
INFO[0001] [最大回應時間: 312.241335ms]
WARN[0001] [失敗請求數量: 0]
=======================================================================
```