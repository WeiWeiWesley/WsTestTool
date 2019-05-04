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
# ./exe/mac_amd64 -H echo.websocket.org

=======================================================================
INFO[0006] [執行數量: 10]
INFO[0006] [執行延遲: 10ns ; 1 nanosecond = 0.0000000001 seconds]
INFO[0006] [成功數量: 10]
INFO[0006] [失敗數量: 0]
INFO[0006] [發送完畢]
=======================================================================
```