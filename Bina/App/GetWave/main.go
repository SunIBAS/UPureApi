package GetWave

import (
	"UPureApi/Bina/App/GetWave/Model"
	"UPureApi/Bina/Utils"
	"UPureApi/Core"
	"UPureApi/Core/DataBase/Bina"
	"UPureApi/Core/DataBase/Bina/Table"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdsMarginedFutures/Market"
	Utils2 "UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/Utils"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"strings"
	"time"
)

type getWaveInfo struct {
	symbols  []string
	updateTs int64
}
type GetWaveKLineInfo struct {
	Interval BinaApis.Interval
	Limit    int64
}
type GetWave struct {
	server           BinaHttpUtils.BinaHttpUtils
	queryParam       Market.KLineListApiParam
	sqlite           Bina.SQLite
	infos            getWaveInfo
	GetWaveKLineInfo GetWaveKLineInfo
}

func NewGetWave(configFile string, getWaveKLineInfo GetWaveKLineInfo) GetWave {
	gw := GetWave{
		server: Utils.CreateServe(configFile),
		queryParam: Market.KLineListApiParam{
			Symbol:    "",
			Interval:  BinaApis.Interval15m,
			StartTime: 0,
			EndTime:   0,
			Limit:     4 * 48, // 2 天
		},
		infos: getWaveInfo{
			symbols: []string{},
		},
		GetWaveKLineInfo: getWaveKLineInfo,
	}
	return gw
}

// Query
// 获取指定 币 在一定长时间内的波动情况
// coinNames = ['SOL', 'DOGE']
func (getWave *GetWave) Query(coinNames []string) map[string][]Utils2.KLine {
	qList := Core.QueryList{
		RetryTimes: 3,
		QueryFunc: func(coinName string) (string, bool) {
			getWave.queryParam.Symbol = fmt.Sprintf("%sUSDT", strings.ToUpper(coinName))
			ret, err := getWave.server.Request(Market.CreateKLineApi(getWave.queryParam))
			if err != nil {
				return "", true
			}
			return ret, false
		},
	}
	m := qList.Query(coinNames)
	klineDataMap := make(map[string][]Utils2.KLine, len(coinNames))
	for k, v := range m {
		klineDataMap[k], _ = Utils2.ParseKLineReturn(v)
	}
	return klineDataMap
}

// Init 初始化
func (getWave *GetWave) Init(dbPath string) {
	getWave.sqlite = Bina.SQLite{
		Path: dbPath, // "D:\\all_code\\UPureApi\\config\\Bina.json",
	}
	if err := getWave.sqlite.Init(); err != nil {
		panic(err)
	}
	getWave.init()
	getWave.initBrackets()
}

// UpdateKLineData 从 bina 服务器拉取最新的 kline 数据
func (getWave *GetWave) UpdateKLineData() Model.KLineSortInfoArr {
	getWave.init()
	//progressbar.New(len(getWave.infos.symbols)) // 设置总进度 100
	idx := 0
	retry := 0
	//symbolCount := len(getWave.infos.symbols)
	symbolCount := 10
	kInfoArr := make(Model.KLineSortInfoArr, symbolCount)
	bar := progressbar.DefaultBytes(
		int64(symbolCount),
		"fetch",
	)
	for {
		if idx < symbolCount {
			symbol := getWave.infos.symbols[idx]
			api := Market.CreateKLineApi(
				Market.KLineListApiParam{
					Symbol:    symbol,
					Interval:  getWave.GetWaveKLineInfo.Interval,
					StartTime: 0,
					EndTime:   0,
					Limit:     getWave.GetWaveKLineInfo.Limit,
				},
			)
			ret, _ := getWave.server.Request(api)
			klines, err := Utils2.ParseKLineReturn(ret)
			if err != nil {
				retry++
				if retry > 3 {
					idx++
					fmt.Println(fmt.Sprintf("[%s] fetch fail", symbol))
				}
				continue
			}
			for _, k := range klines {
				getWave.sqlite.GetDb().Save(Table.ApiKLine2KLine(symbol, k))
			}
			kInfoArr[idx] = Model.BuildKLineSortInfo(symbol, klines)
			bar.Add(1)
			bar.Describe(fmt.Sprintf("[%s]", symbol))
			idx++
			retry = 0
			time.Sleep(time.Millisecond * 50)
		} else {
			break
		}
	}
	return kInfoArr
}

// 排序从 bina 服务器拉取的数据
// 这里是否需要考虑什么算法？
func (getWave *GetWave) sortKLines(klines []Utils2.KLine) {
}

// 初始化，（跟新 交易规则、同时更新最后一次查询时间）
func (getWave *GetWave) init() {
	// 从数据库获取最后一次更新的时间
	eiConfigLastTime := &Table.Config{}
	getWave.sqlite.GetDb().First(&eiConfigLastTime, Table.ExchangeInfoSymbolLastUpdate)
	nowTs := time.Now().UnixMilli()
	lastTime := nowTs + Table.ExchangeInfoSymbolLiveTime
	if eiConfigLastTime.ContentInt64 == 0 || eiConfigLastTime.ContentInt64 < nowTs {
		getWave.initExchangeInfo(lastTime)
		// 保存时添加存活时间，这样子可以更简单进行更新
		eiConfigLastTime.Id = Table.ExchangeInfoSymbolLastUpdate.Id
		eiConfigLastTime.Name = Table.ExchangeInfoSymbolLastUpdate.Name
		eiConfigLastTime.ContentInt64 = nowTs + Table.ExchangeInfoSymbolLiveTime
		getWave.sqlite.GetDb().Save(eiConfigLastTime)
	} else {
		lastTime = eiConfigLastTime.ContentInt64
		fmt.Println(fmt.Sprintf("next time %d ; aftet %s \n", eiConfigLastTime.ContentInt64, Core.MillisecondsToTime(eiConfigLastTime.ContentInt64-nowTs)))
		if getWave.infos.updateTs < nowTs {
			getWave.fetchKLineInfo(lastTime)
		}
	}

}

// 从 bina 服务器拉去最新的交易规则信息
func (getWave *GetWave) initExchangeInfo(lastTime int64) {
	fmt.Println("update exchange info")
	// 获取所有的交易对的交易规则
	api := Market.CreateExchangeInfoApi()
	ret, _ := getWave.server.Request(api)
	info := Utils2.ParseExchangeInfo(ret)
	getWave.infos.symbols = make([]string, len(info.Symbols))
	for idx, symbol := range info.Symbols {
		eiSymbol := Table.ExchangeInfoSymbol{
			BaseAsset:             symbol.BaseAsset,
			BaseAssetPrecision:    symbol.BaseAssetPrecision,
			ContractType:          symbol.ContractType,
			DeliveryDate:          symbol.DeliveryDate,
			LiquidationFee:        Core.StringToFloat64(symbol.LiquidationFee),
			MaintMarginPercent:    Core.StringToFloat64(symbol.MaintMarginPercent),
			MarginAsset:           symbol.MarginAsset,
			MarketTakeBound:       Core.StringToFloat64(symbol.MarketTakeBound),
			MaxMoveOrderLimit:     symbol.MaxMoveOrderLimit,
			OnboardDate:           symbol.OnboardDate,
			Pair:                  symbol.Pair,
			PricePrecision:        symbol.PricePrecision,
			QuantityPrecision:     symbol.QuantityPrecision,
			QuoteAsset:            symbol.QuoteAsset,
			QuotePrecision:        symbol.QuotePrecision,
			RequiredMarginPercent: Core.StringToFloat64(symbol.RequiredMarginPercent),
			Status:                symbol.Status,
			Symbol:                symbol.Symbol,
			TriggerProtect:        Core.StringToFloat64(symbol.TriggerProtect),
			UnderlyingType:        symbol.UnderlyingType,
		}
		getWave.sqlite.GetDb().Save(&eiSymbol)
		symbolDataConfig := Table.Config{
			Id:             symbol.Symbol,
			Name:           Table.KLineDataSymbols.Name,
			ContentString:  symbol.Pair,
			ContentInt64:   lastTime,
			ContentFloat64: 0,
		}
		getWave.infos.symbols[idx] = symbol.Symbol
		getWave.sqlite.GetDb().Save(symbolDataConfig)
	}
	getWave.infos.updateTs = lastTime
}

// 从 数据库 拉去交易对
func (getWave *GetWave) loadSymbolFormDb(lastTime int64) {
	var KLineSymbols []Table.Config
	getWave.sqlite.GetDb().Where(&Table.Config{Name: Table.KLineDataSymbols.Name, ContentInt64: lastTime}).Find(&KLineSymbols)
	getWave.infos.symbols = make([]string, len(KLineSymbols))
	for idx, kls := range KLineSymbols {
		getWave.infos.symbols[idx] = kls.Id
	}
	return
}

// 检查是否需要从新拉去交易对信息
func (getWave *GetWave) fetchKLineInfo(lastTime int64) {
	fmt.Println("update kline info")
	// 一般交易对不可能这么少
	if len(getWave.infos.symbols) < 10 {
		getWave.loadSymbolFormDb(lastTime)
		getWave.infos.updateTs = lastTime
	}
	//api := Market.CreateKLineApi()
}

func (getWave *GetWave) GetExchangeInfoFromDb(symbol string) (info Table.ExchangeInfoSymbol) {
	getWave.sqlite.GetDb().Where(&Table.ExchangeInfoSymbol{Symbol: symbol}).Find(&info)
	return
}

// 从服务器拉取 brackets （可开倍率） 数据
func (getWave *GetWave) initBrackets() {
	var bc Table.Config
	getWave.sqlite.GetDb().Where(&Table.BracketsConfig).Find(&bc)
	if bc.ContentString == Table.BracketsConfigSyncComplete {
		// 已经拉取过
		return
	} else {
		// 重新拉取
		//api := Market.CreateBracketsParamApi(Market.BracketsParam{})
		//ret, _ := getWave.server.Post(api.Path, HttpUtils.DefaultHeader, "")
		ret, _ := ioutil.ReadFile("D:\\all_code\\UPureApi\\config\\brackets.json")
		brackets := Table.ParseBracketsResponseObject(string(ret))
		for _, b := range brackets {
			getWave.sqlite.GetDb().Save(b)
		}
		getWave.sqlite.GetDb().Save(Table.Config{
			Id:             Table.BracketsConfig.Id,
			Name:           Table.BracketsConfig.Name,
			ContentString:  Table.BracketsConfigSyncComplete,
			ContentInt64:   0,
			ContentFloat64: 0,
		})
	}
}
