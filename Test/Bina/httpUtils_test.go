package Bina

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Account/UMFAccountAccount"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Account/UMFAccountBalance"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeLeverage"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeMarginType"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOpen"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrder"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrderAll"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrderDelete"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func runApi(api BinaHttpUtils.Api) {
	server := createServe()
	if ret, err := server.Request(api); err == nil {
		fmt.Println(ret)
	} else {
		panic(err)
	}
}

func TestKLineList(t *testing.T) {
	api := UMFMarketKLine.CreateKLineApi(
		UMFMarketKLine.KLineListApiParam{
			Symbol:    "BTCUSDT",
			Interval:  BinaApis.Interval1m,
			StartTime: 0,
			EndTime:   0,
			Limit:     0,
		},
	)
	runApi(api)
}

func TestOrderAll(t *testing.T) {
	orderQueryApi := UMFTradeOrderAll.CreateOrderAllApi(
		UMFTradeOrderAll.OrderAllParam{
			Symbol: "DOGEUSDT",
		},
	)
	runApi(orderQueryApi)
}

func TestOrderQuery(t *testing.T) {
	orderQueryApi := UMFTradeOpen.CreateOpenApi(
		UMFTradeOpen.OpenParam{
			Symbol: "DOGEUSDT",
		},
	)
	runApi(orderQueryApi)
}

// 测试挂单
func TestTrade(t *testing.T) {
	testParams := []struct {
		qty   float64
		price float64
		order UMFTradeOrder.StopOrder
	}{
		{
			qty:   0.02,
			price: 1000,
			order: UMFTradeOrder.StartOrderLongLimit,
		},
		{
			qty:   0.01,
			price: 5000,
			order: UMFTradeOrder.StartOrderShortLimit,
		},
		{
			qty:   0.02,
			price: 1000,
			order: UMFTradeOrder.StartOrderLongMarket,
		},
		{
			qty:   0.02,
			price: 5000,
			order: UMFTradeOrder.StartOrderShortMarket,
		},
	}
	for _, param := range testParams {
		orderParam := UMFTradeOrder.OrderParam{
			Symbol: "ETHUSDT",
		}
		orderParam.OrderParamOrder(param.qty, param.order, param.price)
		api := UMFTradeOrder.CreateOrderApi(
			orderParam,
		)
		runApi(api)
	}
}

// 测试止盈止损
func TestTradeProfileStop(t *testing.T) {
	testParams := []struct {
		price float64
		qty   float64
		order UMFTradeOrder.StopOrder
		desc  string
	}{
		{
			desc:  "止盈多单",
			price: 3000,
			qty:   1,
			order: UMFTradeOrder.StopOrderLongProfit,
		},
		{
			desc:  "止损多单",
			price: 1000,
			qty:   1,
			order: UMFTradeOrder.StopOrderLongLoss,
		},
		{
			desc:  "止盈空单",
			price: 1000,
			qty:   1,
			order: UMFTradeOrder.StopOrderShortProfit,
		},
		{
			desc:  "止损空单",
			price: 3000,
			qty:   1,
			order: UMFTradeOrder.StopOrderShortLoss,
		},
	}
	for _, param := range testParams {
		fmt.Println(fmt.Sprintf("=====[%s]=====", param.desc))
		stopLongProfitParam := UMFTradeOrder.OrderParam{
			Symbol:          "ETHUSDT",
			Type:            "",
			Side:            "",
			PositionSide:    "",
			ReduceOnly:      "",
			Quantity:        0,
			Price:           0,
			StopPrice:       0,
			ClosePosition:   "",
			ActivationPrice: 0,
			CallbackRate:    0,
			PriceProtect:    "",
		}
		// 止盈多单
		stopLongProfitParam.OrderParamStopProfit(param.price, param.qty, param.order)
		stopLongProfit := UMFTradeOrder.CreateOrderApi(stopLongProfitParam)
		stopLongProfit.Path = "/fapi/v1/order"
		runApi(stopLongProfit)
	}
}

// 设置杠杆和持仓方式
func TestSetLeverageAndMarginType(t *testing.T) {
	setLeverage := UMFTradeLeverage.CreateLeverAgeApi(
		UMFTradeLeverage.LeverAgeParam{
			Symbol:   "ETHUSDT",
			Leverage: 100,
		},
	)
	runApi(setLeverage)
	setMarginType := UMFTradeMarginType.CreateMarginTypeApi(
		UMFTradeMarginType.MarginTypeParam{
			Symbol:     "BTCUSDT",
			MarginType: UMFTradeMarginType.MarginTypeCROSSED,
		},
	)
	runApi(setMarginType)
}

func TestTimestamp(t *testing.T) {
	tt := strconv.FormatInt(time.Now().UnixMilli(), 10)
	fmt.Println(tt)
}

func TestBalance(t *testing.T) {
	server := createServe()
	api := UMFAccountBalance.CreateBalanceApi(
		UMFAccountBalance.BalanceParam{},
	)
	ret, _ := server.Request(api)
	balance := UMFAccountBalance.ParseResponseToBalance(ret)
	fmt.Println(balance)
}

func TestAccount(t *testing.T) {
	server := createServe()
	api := UMFAccountAccount.CreateAccountApi(
		UMFAccountAccount.AccountParam{},
	)
	ret, _ := server.Request(api)
	balance := UMFAccountAccount.ParseResponseToBalance(ret)
	fmt.Println(balance)
}

// 查看挂单
func TestOpenOrder(t *testing.T) {
	server := createServe()
	api := UMFTradeOpen.CreateOpenApi(
		UMFTradeOpen.OpenParam{
			Symbol: "ETHUSDT",
		},
	)
	ret, _ := server.Request(api)
	order := UMFTradeOpen.ParseResponseToBalance(ret)
	fmt.Println(order)
}

// 测试 挂单 之后 撤单
func TestOrderThenReverse(t *testing.T) {
	symbol := "BTCUSDT"
	server := createServe()
	orderParam := UMFTradeOrder.OrderParam{
		Symbol: symbol,
	}
	// 止盈单, {8w} 时 {1} 张 [symbol] {止盈}
	orderParam.OrderParamStopProfit(80000, 1, UMFTradeOrder.StopOrderLongProfit)
	api := UMFTradeOrder.CreateOrderApi(orderParam)
	l, _ := server.RequestL(api, true)
	fmt.Println(l)

	// 等待 2 秒
	time.Sleep(time.Second * 2)

	openRet, _ := server.RequestL(UMFTradeOpen.CreateOpenApi(UMFTradeOpen.OpenParam{
		Symbol: symbol,
	}), true)
	open := UMFTradeOpen.ParseResponseToBalance(openRet)
	// 打印出挂了哪些单
	fmt.Println(open)

	for _, o := range open {
		// 避免意外，这里只 撤销 symbol 的挂单
		if o.Symbol == symbol {
			api := UMFTradeOrderDelete.CreateOrderDeleteApi(UMFTradeOrderDelete.OrderDeleteParam{
				Symbol:  symbol,
				OrderId: o.OrderId,
			})
			deleteRet, _ := server.RequestL(api, true)
			fmt.Println(deleteRet)
		} else {
			fmt.Printf("[Question] %T\r\n", o)
		}
	}
}
