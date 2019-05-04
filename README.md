WsTestTool
===

# Usage

### Start with go: go run main.go -h
### Start with os version: ./exe/${OS_VERSION} -h  
  
```bash
# ./exe/mac_amd64 -h


  -H string
    	Host.
  -P string
    	URL path.
  -d int
    	Time duration between each request. (default 10)
  -h	Usage.
  -n int
    	Test times. (default 1)
  -timeout int
    	Test times. (default 10)
  -w	Watch each resposnes.
```

# Simple test

```bash=
# ./exe/mac_amd64 -H echo.websocket.org -w -n 10

INFO[0000] [Sent count: 1 Response: {"command":"ping"} Time: 203.98563ms]
INFO[0001] [Sent count: 2 Response: {"command":"ping"} Time: 433.792752ms]
INFO[0002] [Sent count: 3 Response: {"command":"ping"} Time: 303.88326ms]
INFO[0002] [Sent count: 4 Response: {"command":"ping"} Time: 224.968427ms]
INFO[0003] [Sent count: 5 Response: {"command":"ping"} Time: 315.91325ms]
INFO[0004] [Sent count: 6 Response: {"command":"ping"} Time: 218.365575ms]
INFO[0004] [Sent count: 7 Response: {"command":"ping"} Time: 1.178186635s]
INFO[0004] [Sent count: 8 Response: {"command":"ping"} Time: 220.85474ms]
INFO[0004] [Sent count: 9 Response: {"command":"ping"} Time: 210.389607ms]
INFO[0005] [Sent count: 10 Response: {"command":"ping"} Time: 214.647233ms]

=======================================================================
INFO[0005] [執行數量: 10]
INFO[0005] [執行延遲: 10ns ; 1 nanosecond = 0.0000000001 seconds]
INFO[0005] [成功數量: 10]
INFO[0005] [總執行時間: 3.524987109s]
INFO[0005] [平均回應時間: 352.49871ms]
INFO[0005] [最大回應時間: 1.178186635s]
INFO[0005] [發送完畢]
=======================================================================
```