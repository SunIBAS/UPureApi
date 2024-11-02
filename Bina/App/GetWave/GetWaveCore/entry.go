package GetWaveCore

import (
	"UPureApi/Bina/App/GetWave/Model"
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/CoreUtils/loggerUtils"
	"UPureApi/Core/DataBase/Bina"
	"UPureApi/Core/DataBase/Bina/Table"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketBrackets"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketExchangeInfo"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeLeverage"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeMarginType"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"path"
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
	queryParam       UMFMarketKLine.KLineListApiParam
	sqlite           Bina.SQLite
	infos            getWaveInfo
	GetWaveKLineInfo GetWaveKLineInfo
	logger           loggerUtils.KeepLogger
	MaxMargin        int
}

// 需要存在下面订单方式才能用来交易(这里使用市价来完成交易)
var orderTypes = []string{
	"MARKET",
	"STOP_MARKET",
	"TAKE_PROFIT_MARKET",
}

type Option func(*GetWave)

func NewGetWave(config BinaHttpUtils.BinaHttpUtilsConfig, getWaveKLineInfo GetWaveKLineInfo, opt ...Option) GetWave {
	gw := GetWave{
		server: BinaHttpUtils.NewBinaHttpUtilsFromConfig(config),
		queryParam: UMFMarketKLine.KLineListApiParam{
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
		logger:           loggerUtils.CreateLogger(path.Join(path.Dir(os.Args[0]), "log.log")),
		MaxMargin:        -1,
	}
	for _, o := range opt {
		o(&gw)
	}
	return gw
}
func (getWave *GetWave) GetLogger() loggerUtils.KeepLogger {
	return getWave.logger
}

// Query
// 获取指定 币 在一定长时间内的波动情况
// coinNames = ['SOL', 'DOGE']
func (getWave *GetWave) Query(coinNames []string) map[string][]UMFMarketKLine.KLineResponse {
	qList := CoreUtils.QueryList{
		RetryTimes: 3,
		QueryFunc: func(coinName string) (string, bool) {
			getWave.queryParam.Symbol = fmt.Sprintf("%sUSDT", strings.ToUpper(coinName))
			ret, err := getWave.server.Request(UMFMarketKLine.CreateKLineApi(getWave.queryParam))
			if err != nil {
				return "", true
			}
			return ret, false
		},
	}
	m := qList.Query(coinNames)
	klineDataMap := make(map[string][]UMFMarketKLine.KLineResponse, len(coinNames))
	for k, v := range m {
		klineDataMap[k] = UMFMarketKLine.ParseKLineResponse(v)
	}
	return klineDataMap
}

// Init 初始化
func (getWave *GetWave) Init(dbPath string) {
	getWave.logger.Info("初始化...")
	getWave.sqlite = Bina.SQLite{
		Path: dbPath, // "D:\\all_code\\UPureApi\\config\\Bina.json",
	}
	if err := getWave.sqlite.Init(); err != nil {
		getWave.logger.Warn("初始化 数据库 失败", "\r\nError:\r\n", err.Error())
		panic(err)
	}
	getWave.logger.Info("初始化 数据库 完成")
	getWave.initPairInfos()
	getWave.initBrackets()
	getWave.logger.Info("初始化 彻底完成")
}

// UpdateKLineData 从 bina 服务器拉取最新的 kline 数据
func (getWave *GetWave) UpdateKLineData() Model.KLineSortInfoArr {
	getWave.logger.Info("拉取最新 K 线数据")
	getWave.logger.Info("先检查是否有交易对数据")
	getWave.initPairInfos()
	//progressbar.New(len(getWave.infos.symbols)) // 设置总进度 100
	idx := 0
	retry := 0
	symbolCount := len(getWave.infos.symbols)
	//symbolCount := 10
	kInfoArr := make(Model.KLineSortInfoArr, symbolCount)
	bar := progressbar.DefaultBytes(
		int64(symbolCount),
		"fetch",
	)
	getWave.logger.Info(fmt.Sprintf("开始拉取 K 线数据，共需要拉取 [%d] 个交易对数据", symbolCount))
	for {
		if idx < symbolCount {
			symbol := getWave.infos.symbols[idx]
			if symbol == "" {
				// 跳过空的交易对
				idx++
				continue
			}
			api := UMFMarketKLine.CreateKLineApi(
				UMFMarketKLine.KLineListApiParam{
					Symbol:    symbol,
					Interval:  getWave.GetWaveKLineInfo.Interval,
					StartTime: 0,
					EndTime:   0,
					Limit:     getWave.GetWaveKLineInfo.Limit,
				},
			)
			ret, _ := getWave.server.Request(api)
			lines := UMFMarketKLine.ParseKLineResponse(ret)
			if len(lines) == 0 {
				retry++
				if retry > 3 {
					idx++
					fmt.Println(fmt.Sprintf("[%s] fetch fail", symbol))
				}
				continue
			}
			// 感觉没必要将 kLine 保存到数据库
			//for _, k := range lines {
			//	getWave.sqlite.GetDb().Save(Table.ApiKLine2KLine(symbol, k))
			//}
			kInfoArr[idx] = Model.BuildKLineSortInfo(symbol, lines)
			bar.Add(1)
			bar.Describe(fmt.Sprintf("[%s]", symbol))
			idx++
			retry = 0
			time.Sleep(time.Millisecond * 50)
		} else {
			break
		}
	}
	getWave.logger.Info(fmt.Sprintf("拉取 K 线数据完成"))
	return kInfoArr
}

// 排序从 bina 服务器拉取的数据
// 这里是否需要考虑什么算法？
//func (getWave *GetWave) sortKLines(klines UMFMarketKLine.KLineResponse) {
//}

// 初始化交易对信息，（跟新 交易规则、同时更新最后一次查询时间）
func (getWave *GetWave) initPairInfos() {
	// 从数据库获取最后一次更新的时间
	getWave.logger.Info("初始化 拉取 交易对数据")
	eiConfigLastTime := &Table.Config{}
	getWave.sqlite.GetDb().First(&eiConfigLastTime, Table.ExchangeInfoSymbolLastUpdate)
	nowTs := time.Now().UnixMilli()
	lastTime := nowTs + Table.ExchangeInfoSymbolLiveTime
	getWave.logger.Info("初始化 拉取 交易对数据")
	if eiConfigLastTime.ContentInt64 == 0 ||
		eiConfigLastTime.ContentInt64 < nowTs {
		getWave.logger.Info("初始化 交易对数据在本地过时需要重新从官网拉取")
		getWave.initExchangeInfo(lastTime)
		// 保存时添加存活时间，这样子可以更简单进行更新
		eiConfigLastTime.Id = Table.ExchangeInfoSymbolLastUpdate.Id
		eiConfigLastTime.Name = Table.ExchangeInfoSymbolLastUpdate.Name
		eiConfigLastTime.ContentInt64 = nowTs + Table.ExchangeInfoSymbolLiveTime
		getWave.sqlite.GetDb().Save(eiConfigLastTime)
	} else {
		getWave.logger.Info("初始化 交易对数据 在数据库拉取")
		lastTime = eiConfigLastTime.ContentInt64
		fmt.Println(fmt.Sprintf("next time %d ; aftet %s \n", eiConfigLastTime.ContentInt64, CoreUtils.MillisecondsToTime(eiConfigLastTime.ContentInt64-nowTs)))
		if getWave.infos.updateTs < nowTs {
			getWave.fetchKLineInfo(lastTime)
		}
	}

}

// 从 bina 服务器拉去最新的交易规则信息
func (getWave *GetWave) initExchangeInfo(lastTime int64) {
	getWave.logger.Info("初始化 正在从官网拉取交易对交易规则数据")
	getWave.logger.Info("交易规则包括，可开最大最小金额/数量，交易金额/数量的最小单位，支持止盈止损规则等")
	getWave.logger.Info("price 金额，qty 数量")
	// 获取所有的交易对的交易规则
	api := UMFMarketExchangeInfo.CreateExchangeInfoApi(UMFMarketExchangeInfo.ExchangeInfoParam{})
	ret, _ := getWave.server.Request(api)
	info := UMFMarketExchangeInfo.ParseExchangeInfoResponse(ret)
	getWave.infos.symbols = make([]string, len(info.Symbols))
	getWave.logger.Info(fmt.Sprintf("初始化 拉取成功，共有 [%d] 个交易对", len(info.Symbols)))
	for idx, symbol := range info.Symbols {
		if symbol.CanOrder(orderTypes) {
			eiSymbol := UMFMarketExchangeInfo.ParseExchangeInfoResponseTable(symbol)
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
			getWave.logger.Info(fmt.Sprintf("交易对 [%s] 交易规则如下", symbol.Pair))
			getWave.logger.Info(fmt.Sprintf("交易价 [最大：%f,最小：%f, 单位：%f]", eiSymbol.FMaxPrice, eiSymbol.FMinPrice, eiSymbol.FPriceStep))
			getWave.logger.Info(fmt.Sprintf("交易量 [最大：%f,最小：%f, 单位：%f]", eiSymbol.FMaxQty, eiSymbol.FMinQty, eiSymbol.FQtyStep))
			getWave.logger.Info(fmt.Sprintf("交易额 [最小：%f]", eiSymbol.FMinNotional))
		} else {
			getWave.logger.Info(fmt.Sprintf("交易对 [%s] 不支持市价止盈止损操作，被踢出 ", symbol.Pair))
		}
	}
	getWave.logger.Info("初始化 拉取记录已被同步到数据库")
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

// 从服务器拉取 brackets （可开倍率） 数据
func (getWave *GetWave) initBrackets() {
	getWave.logger.Info("初始化 开始拉取 杠杆")
	var bc Table.Config
	getWave.sqlite.GetDb().Where(&Table.BracketsConfig).Find(&bc)
	if bc.ContentString == Table.BracketsConfigSyncComplete {
		getWave.logger.Info("初始化 已经在数据库中获取到数据，可以按需从数据库获取")
		// 已经拉取过
		return
	} else {
		getWave.logger.Info("初始化 数据不存在，开始从官网拉取")
		// 重新拉取
		api := UMFMarketBrackets.CreateBracketsApi(UMFMarketBrackets.BracketsParam{})
		//ret, _ := getWave.server.Post(api.Path, HttpUtils.DefaultHeader, "")
		ret, _ := getWave.server.Request(api)
		brackets := UMFMarketBrackets.PParseBracketResponseTable(ret)
		getWave.logger.Info(fmt.Sprintf("初始化 拱拉取到 [%d] 个交易对信息", len(brackets)))
		for _, b := range brackets {
			getWave.sqlite.GetDb().Save(b)
			getWave.logger.Info("[%s] 最大杠杆倍数 [%d]", b.Symbol, b.MaxOpenPosLeverage)
		}
		getWave.logger.Info("初始化 同步到数据库完成")
		getWave.sqlite.GetDb().Save(Table.Config{
			Id:             Table.BracketsConfig.Id,
			Name:           Table.BracketsConfig.Name,
			ContentString:  Table.BracketsConfigSyncComplete,
			ContentInt64:   0,
			ContentFloat64: 0,
		})
	}
}

// GetExchangeInfoFromDb 从 db 获取某个交易对的交易规则
func (getWave *GetWave) GetExchangeInfoFromDb(symbol string) (info Table.ExchangeInfoSymbol) {
	getWave.sqlite.GetDb().Where(&Table.ExchangeInfoSymbol{Symbol: symbol}).Find(&info)
	return
}

// GetBracketsFromDb 获取交易对的交易杠杆限制
func (getWave *GetWave) GetBracketsFromDb(symbol string) (brackets Table.Brackets) {
	getWave.sqlite.GetDb().Where(&Table.Brackets{Symbol: symbol}).Find(&brackets)
	return
}

// SetCrossAndMaxLeverage 设置交易对的最大开仓杠杆和全仓模式
// pair = BTCUSDT
func (getWave *GetWave) SetCrossAndMaxLeverage(pair string) int {
	marginApi := UMFTradeMarginType.CreateMarginTypeApi(UMFTradeMarginType.MarginTypeParam{
		Symbol:     pair,
		MarginType: UMFTradeMarginType.MarginTypeCROSSED,
	})
	marginRet, _ := getWave.server.Request(marginApi)
	fmt.Println(fmt.Sprintf("[margin response] %s\r\n", marginRet))
	bracket := getWave.GetBracketsFromDb(pair)
	leverage := bracket.MaxOpenPosLeverage
	if getWave.MaxMargin > 0 {
		leverage = min(leverage, getWave.MaxMargin)
	}
	if leverage > 1 {
		leverageApi := UMFTradeLeverage.CreateLeverAgeApi(UMFTradeLeverage.LeverAgeParam{
			Symbol:   pair,
			Leverage: int64(leverage),
		})
		leverageRet, _ := getWave.server.Request(leverageApi)
		fmt.Println(fmt.Sprintf("[leverage response] %s\r\n", leverageRet))
	}
	return leverage
}
