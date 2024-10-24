package App

import (
	"UPureApi/Bina/App/GetWave"
	"UPureApi/Bina/Utils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdsMarginedFutures/Market"
	Utils2 "UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/Utils"
	"fmt"
	"testing"
)

var configFile = "D:\\all_code\\UPureApi\\config\\Bina.json"

func TestGetWave(t *testing.T) {
	gw := GetWave.NewGetWave(configFile, GetWave.GetWaveKLineInfo{
		Interval: BinaApis.Interval15m,
		Limit:    4 * 24 * 2,
	})

	dbPath := "D:\\all_code\\UPureApi\\config\\gw.db"
	gw.Init(dbPath)
	gw.UpdateKLineData()
	fmt.Println(gw)

	//ret := gw.Query([]string{
	//	"SOL",
	//	"BTC",
	//	"ETH",
	//})
	//bts, _ := json.Marshal(ret)
	//fmt.Println(string(bts))
	//ioutil.WriteFile("D:\\all_code\\UPureApi\\Test\\Bina\\App\\GetWave_test.json", bts, fs.ModeAppend)
}

func TestPremiumIndex(t *testing.T) {
	api := Market.CreatePremiumIndexApi(
		Market.PremiumIndex{},
	)
	server := Utils.CreateServe(configFile)
	ret, _ := server.Request(api)
	pre := Utils2.ParsePremiumIndexReturn(ret)
	fmt.Println(pre)
}

func TestExchangeInfo(t *testing.T) {
	api := Market.CreateExchangeInfoApi()
	server := Utils.CreateServe(configFile)
	ret, _ := server.Request(api)
	info := Utils2.ParseExchangeInfo(ret)
	fmt.Println(info)
}
