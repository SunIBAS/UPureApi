package UMFTradeOpen

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
)

/* 获取所有历史订单 */

// OpenParam https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Current-All-Open-Orders
type OpenParam struct {
	Symbol     string
	RecvWindow int64
	//Timestamp         int64
}

func (openOrderParams OpenParam) ToMap() BinaHttpUtils.ParamMap {
	m := map[string]string{
		"symbol": BinaApis.CheckEmptyString(openOrderParams.Symbol),
	}
	return m
}

func CreateOpenOrderApi(openOrderParams OpenParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/openOrders",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: openOrderParams,
		BodyParams:  BinaApis.EmptyParams{},
		Sign:        true,
		Header:      HttpUtilsCore.DefaultHeader,
	}
}
