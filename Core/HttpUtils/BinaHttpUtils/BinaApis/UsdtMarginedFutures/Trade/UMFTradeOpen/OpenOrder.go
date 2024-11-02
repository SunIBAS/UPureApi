package UMFTradeOpen

import (
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"encoding/json"
	"fmt"
)

/* 获取所有历史订单 */

// OpenParam https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Current-All-Open-Orders
type OpenParam struct {
	Symbol     string
	RecvWindow int64
	//Timestamp         int64
}

func (openParams OpenParam) ToMap() BinaHttpUtils.ParamMap {
	m := CoreUtils.ApiParamMap{}
	m.SetNotEmptyString("symbol", openParams.Symbol)
	return BinaHttpUtils.ParamMap(m)
}

func CreateOpenApi(openParams OpenParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/openOrders",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: openParams,
		BodyParams:  BinaApis.EmptyParams{},
		Sign:        true,
		Header:      HttpUtilsCore.DefaultHeader,
	}
}

type OpenResponse struct {
	AvgPrice                string `json:"avgPrice"`
	ClientOrderId           string `json:"clientOrderId"`
	CumQuote                string `json:"cumQuote"`
	ExecutedQty             string `json:"executedQty"`
	OrderId                 int64  `json:"orderId"`
	OrigQty                 string `json:"origQty"`
	OrigType                string `json:"origType"`
	Price                   string `json:"price"`
	ReduceOnly              bool   `json:"reduceOnly"`
	Side                    string `json:"side"`
	PositionSide            string `json:"positionSide"`
	Status                  string `json:"status"`
	StopPrice               string `json:"stopPrice"`
	ClosePosition           bool   `json:"closePosition"`
	Symbol                  string `json:"symbol"`
	Time                    int64  `json:"time"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	ActivatePrice           string `json:"activatePrice"`
	PriceRate               string `json:"priceRate"`
	UpdateTime              int64  `json:"updateTime"`
	WorkingType             string `json:"workingType"`
	PriceProtect            bool   `json:"priceProtect"`
	PriceMatch              string `json:"priceMatch"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	GoodTillDate            int    `json:"goodTillDate"`
}

func ParseResponseToBalance(str string) []OpenResponse {
	var resp []OpenResponse
	err := json.Unmarshal([]byte(str), &resp)
	if err != nil {
		fmt.Println(err.Error())
		return []OpenResponse{}
	}
	return resp
}
