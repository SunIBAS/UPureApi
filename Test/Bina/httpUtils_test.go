package Bina

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Account/UMFAccountBalance"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeLeverage"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeMarginType"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOpen"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrder"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrderAll"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func createServe() BinaHttpUtils.BinaHttpUtils {
	configFile := "D:\\all_code\\UPureApi\\config\\Bina.json"
	var config BinaHttpUtils.BinaHttpUtilsConfig
	bs, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(bs, &config)
	return BinaHttpUtils.NewBinaHttpUtilsFromConfig(config)
}

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
	orderQueryApi := UMFTradeOpen.CreateOpenOrderApi(
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
