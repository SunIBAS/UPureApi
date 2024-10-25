package UMFTradeOrder

import "UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade"

type StopOrder struct {
	Side         Trade.SideType
	PositionSide Trade.PositionSideType
	Type         Trade.OrderType
}

// StopOrderLongProfit 止盈多单
var StopOrderLongProfit StopOrder = StopOrder{
	Side:         Trade.SideSell,
	PositionSide: Trade.PositionSideLong,
	Type:         Trade.OrderParamTypeTAKE_PROFIT_MARKET,
}

// StopOrderLongLoss 止损多单
var StopOrderLongLoss StopOrder = StopOrder{
	Side:         Trade.SideSell,
	PositionSide: Trade.PositionSideLong,
	Type:         Trade.OrderParamTypeSTOP_MARKET,
}

// StopOrderShortProfit 止盈空单
var StopOrderShortProfit StopOrder = StopOrder{
	Side:         Trade.SideBuy,
	PositionSide: Trade.PositionSideShort,
	Type:         Trade.OrderParamTypeTAKE_PROFIT_MARKET,
}

// StopOrderShortLoss 止损空单
var StopOrderShortLoss StopOrder = StopOrder{
	Side:         Trade.SideBuy,
	PositionSide: Trade.PositionSideShort,
	Type:         Trade.OrderParamTypeSTOP_MARKET,
}

var StartOrderLongMarket = StopOrder{
	Side:         Trade.SideBuy,
	PositionSide: Trade.PositionSideLong,
	Type:         Trade.OrderParamTypeMARKET,
}
var StartOrderShortMarket = StopOrder{
	Side:         Trade.SideSell,
	PositionSide: Trade.PositionSideShort,
	Type:         Trade.OrderParamTypeMARKET,
}

var StartOrderLongLimit = StopOrder{
	Side:         Trade.SideBuy,
	PositionSide: Trade.PositionSideLong,
	Type:         Trade.OrderParamTypeLIMIT,
}
var StartOrderShortLimit = StopOrder{
	Side:         Trade.SideSell,
	PositionSide: Trade.PositionSideShort,
	Type:         Trade.OrderParamTypeLIMIT,
}
