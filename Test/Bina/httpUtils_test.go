package Bina

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdsMarginedFutures/Market"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdsMarginedFutures/Trade"
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
	api := Market.CreateKLineApi(
		Market.KLineListApiParam{
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
	orderQueryApi := Trade.CreateOrderAllApi(
		Trade.OrderAllParams{
			Symbol: "DOGEUSDT",
		},
	)
	runApi(orderQueryApi)
}

func TestOrderQuery(t *testing.T) {
	orderQueryApi := Trade.CreateOpenOrderApi(
		Trade.OpenOrderParams{
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
		order Trade.StopOrder
	}{
		{
			qty:   0.02,
			price: 1000,
			order: Trade.StartOrderLongLimit,
		},
		{
			qty:   0.01,
			price: 5000,
			order: Trade.StartOrderShortLimit,
		},
		{
			qty:   0.02,
			price: 1000,
			order: Trade.StartOrderLongMarket,
		},
		{
			qty:   0.02,
			price: 5000,
			order: Trade.StartOrderShortMarket,
		},
	}
	for _, param := range testParams {
		orderParam := Trade.OrderParam{
			Symbol: "ETHUSDT",
		}
		orderParam.OrderParamOrder(param.qty, param.order, param.price)
		api := Trade.CreateOrderApi(
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
		order Trade.StopOrder
		desc  string
	}{
		{
			desc:  "止盈多单",
			price: 3000,
			qty:   1,
			order: Trade.StopOrderLongProfit,
		},
		{
			desc:  "止损多单",
			price: 1000,
			qty:   1,
			order: Trade.StopOrderLongLoss,
		},
		{
			desc:  "止盈空单",
			price: 1000,
			qty:   1,
			order: Trade.StopOrderShortProfit,
		},
		{
			desc:  "止损空单",
			price: 3000,
			qty:   1,
			order: Trade.StopOrderShortLoss,
		},
	}
	for _, param := range testParams {
		fmt.Println(fmt.Sprintf("=====[%s]=====", param.desc))
		stopLongProfitParam := Trade.OrderParam{
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
		stopLongProfit := Trade.CreateOrderApi(stopLongProfitParam)
		stopLongProfit.Path = "/fapi/v1/order"
		runApi(stopLongProfit)
	}
}

// 设置杠杆和持仓方式
func TestSetLeverageAndMarginType(t *testing.T) {
	setLeverage := Trade.CreateLeverAgeParam(
		Trade.LeverAgeParam{
			Symbol:   "ETHUSDT",
			Leverage: 100,
		},
	)
	runApi(setLeverage)
	setMarginType := Trade.CreateMarginTypeParam(
		Trade.MarginTypeParam{
			Symbol:     "BTCUSDT",
			MarginType: Trade.MarginTypeCROSSED,
		},
	)
	runApi(setMarginType)
}

func TestTimestamp(t *testing.T) {
	tt := strconv.FormatInt(time.Now().UnixMilli(), 10)
	fmt.Println(tt)
}
