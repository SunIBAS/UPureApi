package Trade

import (
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
)

/* 获取所有历史订单 */

// OpenOrderParams https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Current-All-Open-Orders
type OpenOrderParams struct {
	Symbol     string
	RecvWindow int64
	//Timestamp         int64
}

func (openOrderParams OpenOrderParams) ToMap() BinaHttpUtils.ParamMap {
	m := map[string]string{
		"symbol": BinaApis.CheckEmptyString(openOrderParams.Symbol),
	}
	return m
}

func CreateOpenOrderApi(openOrderParams OpenOrderParams) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/openOrders",
		HttpMethod:  HttpUtils.HttpGet,
		QueryParams: openOrderParams,
		BodyParams:  BinaApis.EmptyParams{},
		Sign:        true,
		Header:      HttpUtils.DefaultHeader,
	}
}
