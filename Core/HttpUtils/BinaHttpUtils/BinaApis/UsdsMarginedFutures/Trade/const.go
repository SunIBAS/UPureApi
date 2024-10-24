package Trade

type OrderType string

const (
	OrderParamTypeLIMIT  OrderType = "LIMIT"
	OrderParamTypeMARKET OrderType = "MARKET"
	// stop 是止损，逻辑区别是
	// stop 是止损限价，stop 是止损限价
	// 当价格波动较大时，假设止损目标是价格突破 100，
	// 当价格突破 100 后，限价单会以 100 卖出，如果瞬间低于100则交易失败
	// 				   市价单则会无论当前什么价格，都成交，可能成交价格为 99
	// 限价成交价格为 >= 指定价格，价格低于则成交失败
	// 市价成家价格为 ≈≈ 指定价格
	OrderParamTypeSTOP                 OrderType = "STOP"
	OrderParamTypeSTOP_MARKET          OrderType = "STOP_MARKET"
	OrderParamTypeTAKE_PROFIT          OrderType = "TAKE_PROFIT"
	OrderParamTypeTAKE_PROFIT_MARKET   OrderType = "TAKE_PROFIT_MARKET"
	OrderParamTypeTRAILING_STOP_MARKET OrderType = "TRAILING_STOP_MARKET"
)

type SideType string

const (
	SideSell SideType = "SELL"
	SideBuy  SideType = "BUY"
)

type PositionSideType string

const (
	PositionSideLong  PositionSideType = "LONG"
	PositionSideShort PositionSideType = "SHORT"
)

type BoolType string

const (
	BoolTrue  BoolType = "TRUE"
	BoolFalse BoolType = "FALSE"
)
