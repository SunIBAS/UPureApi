# 编写说明

- API 编写

[githu API](https://binance-docs.github.io/apidocs/futures/cn/#trade-3)

[现货 API DOC](https://developers.binance.com/docs/zh-CN/binance-spot-api-docs/rest-api)

[合约 API DOC](https://developers.binance.com/docs/zh-CN/derivatives)

```text
└─UsdsMarginedFutures   U本位合约
    ├─Market             行情接口
    │      List.go
    │
    └─Trade             交易接口
            All.go      查询所有挂单
            const.go    常量
            OpenOrder.go    开单
            Order.go    历史订单
            Query.go    通过特定条件查询订单
            readme.md
```

- 文件命名规则

1. 目录以 文档的 URL 中 rest-api 的前一个单词为主
2. 接口名称可以是 URL 最后一个短语或请求的最后一个 单词
3. 当存在 POST PUT DELETE GET 公用一个 api 时，请使用 2 中的第一个方法

![img.png](../../../image/BinaApiNamed.png)

- API 编写方法

```go
package Trade

import (
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
)

// file : TestOne.go
//type [驼峰是文件名]Params struct
type TestOneParams struct {
	// 这里将所有需要提交服务器的参数列出
	Symbol string
}

func (testOneParams TestOneParams) ToMap() BinaHttpUtils.ParamMap {
	m := map[string]string {
		"symbol": BinaApis.CheckEmptyString(testOneParams.Symbol),
    }
	return m
}
func CreateOrderAllApi(testOneParams TestOneParams) BinaHttpUtils.Api {
	// 这里根据实际写，有的接口可能同时需要 query 和 body 参数，这里只需要 query
	// 所以 bodyParams 用来一个空参数
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/allOrders",
		HttpMethod:  HttpUtils.HttpGet,
		QueryParams: testOneParams,
		BodyParams:  BinaApis.EmptyParams{},
		Sign:        true,
		Header:      HttpUtils.DefaultHeader,
    }
}
```