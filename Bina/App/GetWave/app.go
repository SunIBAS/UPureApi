package GetWave

import (
	"UPureApi/Bina/App/GetWave/GetWaveCore"
	"UPureApi/Bina/App/GetWave/Model"
	"UPureApi/Bina/Utils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type GetWaveAppConfig struct {
	// BinaHttpUtils 的配置
	Proxy struct {
		Porto string `json:"porto"`
		Host  string `json:"host"`
		Port  string `json:"port"`
	} `json:"proxy"`
	Key struct {
		ApiKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
	} `json:"key"`
	BaseUrl string `json:"base_url"`
	Log     bool   `json:"log"`
	// 数据库路径
	DbPath     string  `json:"dbPath"`
	MaxMargins int     `json:"maxMargins"`
	StopRate   float64 `json:"stopRate"`
	ProfitRate float64 `json:"profitRate"`
}

func ParseGetWaveAppConfig(configPath string) GetWaveAppConfig {
	//configFile := "D:\\all_code\\UPureApi\\config\\Bina.json"
	var config GetWaveAppConfig
	bs, _ := ioutil.ReadFile(configPath)
	json.Unmarshal(bs, &config)
	return config
}
func (config GetWaveAppConfig) GetServerConfig() BinaHttpUtils.BinaHttpUtilsConfig {
	return BinaHttpUtils.BinaHttpUtilsConfig{
		Proxy: struct {
			Porto string `json:"porto"`
			Host  string `json:"host"`
			Port  string `json:"port"`
		}{
			Porto: config.Proxy.Porto,
			Host:  config.Proxy.Host,
			Port:  config.Proxy.Port,
		},
		Key: struct {
			ApiKey    string `json:"api_key"`
			SecretKey string `json:"secret_key"`
		}{
			ApiKey:    config.Key.ApiKey,
			SecretKey: config.Key.SecretKey,
		},
		BaseUrl: config.BaseUrl,
		Log:     config.Log,
	}
}

// CreateApp
// var configFile = "D:\\all_code\\UPureApi\\config\\Bina.json"
func CreateApp(configFile string, opt ...GetWaveCore.Option) {
	config := ParseGetWaveAppConfig(configFile)
	if config.MaxMargins > 1 {
		opt = append(opt, func(wave *GetWaveCore.GetWave) {
			wave.MaxMargin = config.MaxMargins
		})
	}
	// 创建 GetWave 实例
	gw := GetWaveCore.NewGetWave(config.GetServerConfig(), GetWaveCore.GetWaveKLineInfo{
		Interval: BinaApis.Interval15m,
		Limit:    4 * 24 * 2,
	}, opt...)

	logger := gw.GetLogger()

	dbPath := config.DbPath
	gw.Init(dbPath)
	//gw.UpdateKLineData()
	fmt.Println(gw)
	buyCoin := GetWaveCore.CreateBuyCoin(gw)
	buyCoin.UpdateRate(config.StopRate, config.ProfitRate)

	for {
		logger.Info("开始进入循环")
		state, err := buyCoin.GetBuyCoinState()
		// 获取错误，可能是比较严重的错误，不能直接继续，会影响开单，导致订单错乱
		if err != nil {
			time.Sleep(time.Second * 60)
			continue
		}
		if state.HasTrading() {
			logger.Info("当前有交易未完成，稍微进行等待")
			time.Sleep(time.Second * 30)
		} else {
			pair := selectOneSymbolPair(
				gw.UpdateKLineData(),
			)
			if pair != "" {
				buyCoin.AddPair(pair)
				// 每完成一次拉取后等待 15 分钟
				time.Sleep(time.Second * 60)
			} else {
				logger.Info("当前没有合适的交易对，稍微进行等待")
				// 如果没有得搞，那要快一点找下一个
				// 每完成一次拉取后等待 一分钟
				time.Sleep(time.Second * 60)
			}
		}
	}
}

func selectOneSymbolPair(arr Model.KLineSortInfoArr) string {
	ret := Utils.FindAndSortVolatileStocks(arr, 10, 6)
	if len(ret) > 0 {
		return ret[0].Name
	}
	return ""
}

// package main
//
//import (
//	"fmt"
//	"math"
//)
//
//type StockData map[string][]float64
//
//func findMostVolatile(stockData StockData) (string, float64) {
//	maxVolatility := 0.0
//	selectedStock := ""
//
//	for stock, prices := range stockData {
//		// 确保有足够的数据点
//		if len(prices) < 2 {
//			continue
//		}
//
//		// 计算各时间段的波动幅度
//		vol15 := math.Abs(1 - prices[0])
//		vol30 := math.Abs(1 - prices[1])
//		vol45 := math.Abs(1 - prices[2])
//
//		// 确保波动是持续的 (短时间波动比长期波动大)
//		if vol15 > vol30 && vol30 > vol45 {
//			// 找到波动最大的股票
//			if vol15 > maxVolatility {
//				maxVolatility = vol15
//				selectedStock = stock
//			}
//		}
//	}
//
//	return selectedStock, maxVolatility
//}
//
//func main() {
//	// 假设的股票数据结构
//	data := StockData{
//		"A": {0.99, 0.94, 0.93},
//		"B": {1.02, 0.98, 0.99},
//		"C": {1.10, 1.05, 1.04},
//	}
//
//	stock, volatility := findMostVolatile(data)
//	fmt.Printf("Selected stock: %s with volatility: %.2f%%\n", stock, volatility*100)
//}
