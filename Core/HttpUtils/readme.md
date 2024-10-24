# 说明

- HttpUtils 使用方法

```go
package main

import (
	"UPureApi/Core/HttpUtils"
	"net/http"
	"bytes"
)

func main() {
	// 配置代理
	//{
	//  "Proto": "http",
	//  "Host": "127.0.0.1",
	//  "Port": "7890"
	//}
	opt := HttpUtils.SetProxy(
		"http",
		"127.0.0.1",
		"7890",
	)
	req := HttpUtils.CreateRequest(opt)
	// get 请求
	req.ToRequest(HttpUtils.HttpGet, "", HttpUtils.DefaultHeader , nil)
	// post 请求
	parameters := map[string] string {
		"Name": "IBAS",
		"Age": "123",
    }
	var paramsStr = HttpUtils.Params2string(parameters, HttpUtils.ToString)
	req.ToRequest(HttpUtils.HttpPost, "", HttpUtils.DefaultHeader, bytes.NewBuffer([]byte(paramsStr)))
}
```