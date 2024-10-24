package Trade

import (
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"strconv"
)

// https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Query-Order
type OrderQueryParams struct {
	Symbool           string
	OrderId           int64
	OrigClientOrderId string
	RecvWindow        int64
	//Timestamp         int64
}

func (orderQueryParams OrderQueryParams) ToMap() BinaHttpUtils.ParamMap {
	m := map[string]string{
		"symbol": BinaApis.CheckEmptyString(orderQueryParams.Symbool),
	}
	if orderQueryParams.OrderId != 0 {
		m["orderId"] = strconv.FormatInt(orderQueryParams.OrderId, 10)
	}
	if orderQueryParams.RecvWindow != 0 {
		m["recvWindow"] = strconv.FormatInt(orderQueryParams.RecvWindow, 10)
	}
	if orderQueryParams.OrigClientOrderId != "" {
		m["origClientOrderId"] = orderQueryParams.OrigClientOrderId
	}
	return m
}

func CreateOrderQueryApi(orderQueryParams OrderQueryParams) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/order",
		HttpMethod:  HttpUtils.HttpGet,
		QueryParams: orderQueryParams,
		BodyParams:  BinaApis.EmptyParams{},
		Sign:        true,
		Header:      HttpUtils.DefaultHeader,
	}
}
