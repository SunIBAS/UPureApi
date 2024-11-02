package Bina

import (
	"UPureApi/Bina/App/GetWave"
	"UPureApi/Bina/App/GetWave/GetWaveCore"
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrder"
	"fmt"
	"testing"
)

//func createServe() BinaHttpUtils.BinaHttpUtils {
//	configFile := "D:\\all_code\\UPureApi\\config\\Bina.json"
//	var config BinaHttpUtils.BinaHttpUtilsConfig
//	bs, _ := ioutil.ReadFile(configFile)
//	json.Unmarshal(bs, &config)
//	return BinaHttpUtils.NewBinaHttpUtilsFromConfig(config)
//}

// 下一些莫名其妙的单
func OrderSome() {
	server := createServe()

	symbol := "BTCUSDT"
	orderParam := UMFTradeOrder.OrderParam{
		Symbol: symbol,
	}
	// 止盈单, {8w} 时 {1} 张 [symbol] {止盈}
	orderParam.OrderParamStopProfit(80000, 1, UMFTradeOrder.StopOrderLongProfit)
	api := UMFTradeOrder.CreateOrderApi(orderParam)
	server.RequestL(api, true)
	// 止盈单, {1H} 时 {1} 张 [symbol] {止盈}
	orderParam.OrderParamStopProfit(50000, 1, UMFTradeOrder.StopOrderShortProfit)
	api = UMFTradeOrder.CreateOrderApi(orderParam)
	server.RequestL(api, true)
}

func TestBuyCoin(t *testing.T) {
	config := GetWave.GetWaveAppConfig{
		Proxy: struct {
			Porto string `json:"porto"`
			Host  string `json:"host"`
			Port  string `json:"port"`
		}{
			Porto: "http",
			Host:  "localhost",
			Port:  "7890",
		},
		Key: struct {
			ApiKey    string `json:"api_key"`
			SecretKey string `json:"secret_key"`
		}{
			ApiKey:    "BSNdoFJX9D3qmnOw6iIfKawJuqozWQICQtbplecWo9er8x2xKpgkfQxS4NJBUg2d",
			SecretKey: "RsPHKzK9EN1V5w0xa5go1MQHZiBMIC12pAT3ObwkQ1uVa3ozqdtAPWiKCJa4SlU9",
		},
		BaseUrl: "https://fapi.binance.com",
		Log:     false,
		DbPath:  "D:\\all_code\\UPureApi\\config\\gw.db",
	}
	opts := []GetWaveCore.Option{
		// 设置最大的杠杆倍数
		func(wave *GetWaveCore.GetWave) {
			wave.MaxMargin = 45
		},
	}
	gw := GetWaveCore.NewGetWave(config.GetServerConfig(), GetWaveCore.GetWaveKLineInfo{
		Interval: BinaApis.Interval15m,
		Limit:    4 * 24 * 2,
	}, opts...)
	gw.Init(config.DbPath)
	//OrderSome()

	buyCoin := GetWaveCore.CreateBuyCoin(gw)
	// 更新状态
	buyCoin.SetBuyCoinState()
	buyCoin.AddPair("ETHUSDT")
	fmt.Println(buyCoin)

	for {
	}
}

//// 获取持仓情况
//func TestBuyCoinOrder(t *testing.T) {
//	//server := createServe()
//	//api := UMFTradeOpen.CreateOpenOrderApi()
//}

func TestPrices(t *testing.T) {
	price := 360.65
	fixLen := 2
	stopRate := 0.012
	profitRate := 0.024
	longStop := CoreUtils.Fix(price*(1-stopRate), fixLen)
	fmt.Println(fmt.Sprintf("long stop: %f\t%f", longStop, longStop/price))
	longProfit := CoreUtils.Fix(price*(1+profitRate), fixLen)
	fmt.Println(fmt.Sprintf("long profit: %f\t%f", longProfit, longProfit/price))
	shortStop := CoreUtils.Fix(price*(1+stopRate), fixLen)
	fmt.Println(fmt.Sprintf("short stop: %f\t%f", shortStop, shortStop/price))
	shortProfit := CoreUtils.Fix(price*(1-profitRate), fixLen)
	fmt.Println(fmt.Sprintf("short profit: %f\t%f", shortProfit, shortProfit/price))
}

func TestGetPrice(t *testing.T) {
	server := createServe()
	api := UMFMarketKLine.CreateKLineApi(UMFMarketKLine.KLineListApiParam{
		Symbol:    "BTCUSDT",
		Interval:  BinaApis.Interval1m, // 使用 1m 即可
		StartTime: 0,
		EndTime:   0,
		Limit:     1, // 获取 1 个 k 柱 即可
	})
	kRet, _ := server.Request(api)
	kResp := UMFMarketKLine.ParseKLineResponse(kRet)
	fmt.Println(kResp)
}
