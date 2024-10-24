package Trade

import (
	"UPureApi/Core"
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
)

/* 下单 */

// OrderParam https://binance-docs.github.io/apidocs/futures/cn/#trade-3
// OrderParam https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api
type OrderParam struct {
	Symbol       string
	Type         OrderType
	Side         SideType
	PositionSide PositionSideType
	ReduceOnly   BoolType
	Quantity     float64
	Price        float64
	//newClientOrderId
	StopPrice       float64
	ClosePosition   BoolType
	ActivationPrice float64
	CallbackRate    float64
	TimeInForce     string
	//workingType
	PriceProtect BoolType
	//NewOrderRespType string // ACK RESULT
	//PriceMatch
	//selfTradePreventionMode
	//goodTillDate
	//recvWindow
}

func (orderParam OrderParam) ToMap() BinaHttpUtils.ParamMap {
	m := Core.ApiParamMap{
		"symbol": BinaApis.CheckEmptyString(orderParam.Symbol),
		//"timeInForce": "GTC", // 一年有效
	}
	m.SetNotEmptyString("type", string(orderParam.Type))
	m.SetNotEmptyString("side", string(orderParam.Side))
	m.SetNotEmptyString("positionSide", string(orderParam.PositionSide))
	m.SetNotEmptyString("reduceOnly", string(orderParam.ReduceOnly))
	m.SetNotZeroFloat64String("quantity", orderParam.Quantity)
	m.SetNotZeroFloat64String("price", orderParam.Price)
	m.SetNotZeroDecimal("stopPrice", orderParam.StopPrice)
	m.SetNotEmptyString("closePosition", string(orderParam.ClosePosition))
	m.SetNotZeroFloat64String("activationPrice", orderParam.ActivationPrice)
	m.SetNotZeroFloat64String("callbackRate", orderParam.CallbackRate)
	m.SetNotEmptyString("priceProtect", string(orderParam.PriceProtect))
	m.SetNotEmptyString("timeInForce", string(orderParam.TimeInForce))
	return BinaHttpUtils.ParamMap(m)
}

// OrderParamLimit 限价交易
// quantity price 下单的量和交易价格
func (orderParam *OrderParam) OrderParamLimit(quantity, price float64) {
	orderParam.Type = OrderParamTypeLIMIT
	orderParam.Quantity = quantity
	orderParam.Price = price
	//return orderParam
}

// OrderParamMarket 市价交易
// quantity 下单的量
func (orderParam *OrderParam) OrderParamMarket(quantity float64) {
	orderParam.Type = OrderParamTypeMARKET
	orderParam.Quantity = quantity
	orderParam.Price = 0
}

type StopOrder struct {
	Side         SideType
	PositionSide PositionSideType
	Type         OrderType
}

// StopOrderLongProfit 止盈多单
var StopOrderLongProfit StopOrder = StopOrder{
	Side:         SideSell,
	PositionSide: PositionSideLong,
	Type:         OrderParamTypeTAKE_PROFIT_MARKET,
}

// StopOrderLongLoss 止损多单
var StopOrderLongLoss StopOrder = StopOrder{
	Side:         SideSell,
	PositionSide: PositionSideLong,
	Type:         OrderParamTypeSTOP_MARKET,
}

// StopOrderShortProfit 止盈空单
var StopOrderShortProfit StopOrder = StopOrder{
	Side:         SideBuy,
	PositionSide: PositionSideShort,
	Type:         OrderParamTypeTAKE_PROFIT_MARKET,
}

// StopOrderShortLoss 止损空单
var StopOrderShortLoss StopOrder = StopOrder{
	Side:         SideBuy,
	PositionSide: PositionSideShort,
	Type:         OrderParamTypeSTOP_MARKET,
}

var StartOrderLongMarket = StopOrder{
	Side:         SideBuy,
	PositionSide: PositionSideLong,
	Type:         OrderParamTypeMARKET,
}
var StartOrderShortMarket = StopOrder{
	Side:         SideSell,
	PositionSide: PositionSideShort,
	Type:         OrderParamTypeMARKET,
}

var StartOrderLongLimit = StopOrder{
	Side:         SideBuy,
	PositionSide: PositionSideLong,
	Type:         OrderParamTypeLIMIT,
}
var StartOrderShortLimit = StopOrder{
	Side:         SideSell,
	PositionSide: PositionSideShort,
	Type:         OrderParamTypeLIMIT,
}

// OrderParamStopProfit 止盈止损单
// qty 不指定表示全部成交
func (orderParam *OrderParam) OrderParamStopProfit(price float64, qty float64, order StopOrder) {
	// 不限定 量 则 全出
	if qty == 0 {
		orderParam.ClosePosition = BoolTrue
	} else {
		orderParam.Quantity = qty
	}
	orderParam.StopPrice = price

	orderParam.Side = order.Side
	orderParam.PositionSide = order.PositionSide
	// 推荐使用市价止损
	orderParam.Type = order.Type
}

// OrderParamOrder 开单
func (orderParam *OrderParam) OrderParamOrder(qty float64, order StopOrder, price float64) {
	orderParam.Type = order.Type
	if orderParam.Type == OrderParamTypeLIMIT {
		orderParam.TimeInForce = "GTC"
		orderParam.Price = price
	}
	orderParam.Side = order.Side
	orderParam.PositionSide = order.PositionSide
	orderParam.Quantity = qty
}

func CreateOrderApi(param OrderParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/order",
		HttpMethod:  HttpUtils.HttpPost,
		QueryParams: param,
		BodyParams:  nil,
		Sign:        true,
		Header:      nil,
	}
}
