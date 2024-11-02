package App

import (
	"UPureApi/Bina/App/GetWave"
	"UPureApi/Bina/App/GetWave/GetWaveCore"
	"UPureApi/Bina/App/GetWave/Model"
	"UPureApi/Bina/Utils"
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketExchangeInfo"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketPremiumIndex"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"testing"
)

var configFile = "D:\\all_code\\UPureApi\\config\\Bina.json"
var appConfigFile = "D:\\all_code\\UPureApi\\config\\GetWaveAppConfig.json"

func TestGetWaveApp(t *testing.T) {
	opts := []GetWaveCore.Option{
		// 设置最大的杠杆倍数
		func(wave *GetWaveCore.GetWave) {
			wave.MaxMargin = 45
		},
	}
	GetWave.CreateApp(appConfigFile, opts...)
}

func TestGetWave(t *testing.T) {
	opts := []GetWaveCore.Option{
		// 设置最大的杠杆倍数
		func(wave *GetWaveCore.GetWave) {
			wave.MaxMargin = 45
		},
	}
	gw := GetWaveCore.NewGetWave(Utils.ParseConfig(configFile), GetWaveCore.GetWaveKLineInfo{
		Interval: BinaApis.Interval15m,
		Limit:    4 * 24 * 2,
	}, opts...)

	dbPath := "D:\\all_code\\UPureApi\\config\\gw.db"
	gw.Init(dbPath)
	//gw.UpdateKLineData()
	fmt.Println(gw)
	buyCoin := GetWaveCore.CreateBuyCoin(gw)

	buyCoin.AddCoin("HMSTR")

	for {
		lines := gw.UpdateKLineData()
		bs, _ := json.Marshal(lines)
		ioutil.WriteFile("D:\\all_code\\UPureApi\\config\\klines.json", bs, fs.ModeAppend)
		ret := Utils.FindAndSortVolatileStocks(lines, 2, 1)
		fmt.Println(ret)
	}

	//ret := gw.Query([]string{
	//	"SOL",
	//	"BTC",
	//	"ETH",
	//})
	//bts, _ := json.Marshal(ret)
	//fmt.Println(string(bts))
	//ioutil.WriteFile("D:\\all_code\\UPureApi\\Test\\Bina\\App\\GetWave_test.json", bts, fs.ModeAppend)
}

func TestSortKL(t *testing.T) {
	bs, _ := ioutil.ReadFile("D:\\all_code\\UPureApi\\config\\klines.json")
	var lines Model.KLineSortInfoArr
	json.Unmarshal(bs, &lines)
	ret := Utils.FindAndSortVolatileStocks(lines, 10, 6)
	fmt.Println(ret)
}

func TestPremiumIndex(t *testing.T) {
	api := UMFMarketPremiumIndex.CreatePremiumIndexApi(
		UMFMarketPremiumIndex.PremiumIndexParam{},
	)
	server := Utils.CreateServe(configFile)
	ret, _ := server.Request(api)
	pre := UMFMarketPremiumIndex.ParsePremiumIndexResponse(ret)
	fmt.Println(pre)
}

func TestExchangeInfo(t *testing.T) {
	api := UMFMarketExchangeInfo.CreateExchangeInfoApi(UMFMarketExchangeInfo.ExchangeInfoParam{})
	server := Utils.CreateServe(configFile)
	ret, _ := server.Request(api)
	info := UMFMarketExchangeInfo.ParseExchangeInfoResponse(ret)
	fmt.Println(info)
}

func TestP(t *testing.T) {
	fmt.Println(CoreUtils.GetPointLen(0.12345))
	fmt.Println(CoreUtils.GetPointLen(0.12000))
}
