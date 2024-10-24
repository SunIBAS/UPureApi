package Trade

import (
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"strconv"
)

/* 获取所有历史订单 */

// OrderAllParams https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/All-Orders
type OrderAllParams struct {
	Symbol     string
	OrderId    int64
	StartTime  int64
	EndTime    int64
	Limit      int64
	RecvWindow int64
	//Timestamp         int64
}

var OrderAllLimit = BinaApis.IntRangeParam{
	Max:     1000,
	Min:     10,
	Default: 500,
}

func (orderAllParams OrderAllParams) ToMap() BinaHttpUtils.ParamMap {
	m := map[string]string{
		"symbol": BinaApis.CheckEmptyString(orderAllParams.Symbol),
		"limit":  strconv.FormatInt(OrderAllLimit.Get(orderAllParams.Limit), 10),
	}
	if orderAllParams.StartTime > 1000 {
		m["startTime"] = strconv.FormatInt(orderAllParams.StartTime, 10)
	}
	if orderAllParams.EndTime > 1000 {
		m["endTime"] = strconv.FormatInt(orderAllParams.EndTime, 10)
	}
	if orderAllParams.OrderId > 1000 {
		m["orderId"] = strconv.FormatInt(orderAllParams.OrderId, 10)
	}
	if orderAllParams.RecvWindow > 1000 {
		m["recvWindow"] = strconv.FormatInt(orderAllParams.RecvWindow, 10)
	}
	return m
}

func CreateOrderAllApi(orderAllParams OrderAllParams) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/allOrders",
		HttpMethod:  HttpUtils.HttpGet,
		QueryParams: orderAllParams,
		BodyParams:  BinaApis.EmptyParams{},
		Sign:        true,
		Header:      HttpUtils.DefaultHeader,
	}
}
