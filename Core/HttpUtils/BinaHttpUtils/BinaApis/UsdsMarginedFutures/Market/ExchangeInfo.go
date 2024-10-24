package Market

import (
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
)

type ExchangeInfo struct {
}

func (exchangeInfo ExchangeInfo) ToMap() BinaHttpUtils.ParamMap {
	return map[string]string{
		//"showall": "true",
	}
}

// CreateExchangeInfoApi https://binance-docs.github.io/apidocs/futures/cn/#api
func CreateExchangeInfoApi() BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/exchangeInfo",
		HttpMethod:  HttpUtils.HttpGet,
		QueryParams: ExchangeInfo{},
		BodyParams:  nil,
		Sign:        false,
		Header:      nil,
		NoTimeStamp: true,
	}
}
