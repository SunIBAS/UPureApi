package Market

import (
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
)

/* 获取所有交易对的最大最小交易杠杆 */

// BracketsParam https://www.binance.com/bapi/futures/v1/friendly/future/common/brackets
// https://www.binance.com/zh-CN/futures/trading-rules/perpetual/leverage-margin
type BracketsParam struct{}

func (bracketsParam BracketsParam) ToMap() BinaHttpUtils.ParamMap {
	return map[string]string{}
}

func CreateBracketsParamApi(bracketsParam BracketsParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		NoTimeStamp: false,
		Path:        "https://www.binance.com/bapi/futures/v1/friendly/future/common/brackets",
		HttpMethod:  HttpUtils.HttpPost,
		QueryParams: bracketsParam,
		BodyParams:  nil,
		Sign:        false,
		Header:      nil,
	}
}
