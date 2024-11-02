package main

import (
	"UPureApi/Bina/App/GetWave/GetWaveCore"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"encoding/json"
	"fmt"
	"time"
)

var configFile = "D:\\all_code\\UPureApi\\config\\Bina.json"

func main() {
	// 创建一个 读取 类
	gw := GetWaveCore.NewGetWave(configFile, GetWaveCore.GetWaveKLineInfo{
		Interval: BinaApis.Interval15m,
		Limit:    4 * 24 * 2,
	})

	// 数据库路径
	dbPath := "D:\\all_code\\UPureApi\\config\\gw.db"
	// 初始化（拉取最新的交易规则，主要是交易对信息）
	gw.Init(dbPath)
	info := gw.GetExchangeInfoFromDb("BTCUSDT")
	fmt.Println(info)
	bracket := gw.GetBracketsFromDb("BTCUSDT")
	fmt.Println(bracket)
	for {
		// 获取所有 k 线数据
		kinfo := gw.UpdateKLineData()
		bts, _ := json.Marshal(kinfo)
		fmt.Println(bts)
		// 排序 k 线数据
		time.Sleep(time.Hour)
	}
}
