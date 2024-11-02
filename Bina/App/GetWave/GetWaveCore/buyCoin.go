package GetWaveCore

import (
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/CoreUtils/SingleTask"
	"UPureApi/Core/DataBase/Bina/Table"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Account/UMFAccountAccount"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOpen"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrder"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Trade/UMFTradeOrderDelete"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// 由于这里的流程也比较复杂，因此决定在单独的文件中进行设计

var balanceCoinName = "USDT"

// 因为最小交易金额为 20，这里确保一定可以执行，取 25 * 2 = 50
var baseBalance float64 = 40
var stopRate = 0.011
var profitRate = 0.022

type tradingInfo struct {
	// 买了什么币 (BTC)
	coinName string
	// 买了什么币 (BTCUSDT)
	pair string
	// 什么时候开始买的
	createTime string
	// 预期结束时间
	endTime string
	// 是否存在 多 空 交易
	long  bool
	short bool
	// 当前 多空 收益
	longProfit  float64
	shortProfit float64
	// 当前 多空 张数 (用于平掉)
	longSz  float64
	shortSz float64
	// 做多方向的交易Id
	longTradeId int64
	// 做空方向的交易Id
	shortTradeId int64
	// 做多止盈 Id
	longProfitId int64
	// 做空止损 Id
	longStopId int64
	// 做空止盈 Id
	shortProfitId int64
	// 做空止损 Id
	shortStopId int64
}
type buyCoinState struct {
	// 当前剩余可以购买的 usdt 的余额
	currentUsdt float64
	// 当前的交易情况（是否购买了某一种币，开始时间）
	// 虽然这里是复数，但是我们只假定这里只有一种币的交易
	tradingInfo []tradingInfo
}

func (state buyCoinState) HasTrading() bool {
	have := len(state.tradingInfo) > 0
	if have {
		for _, s := range state.tradingInfo {
			if !have {
				break
			}
			have = have && (s.long || s.short)
		}
	}
	return have
}

type BuyCoin struct {
	// 打算购买的币的列表（加入到这个列表不代表一定会购买）
	coinList []string
	// 当前账户状态
	currentState buyCoinState
	//server       BinaHttpUtils.BinaHttpUtils
	//sqlite       Bina.SQLite
	getWave GetWave

	// 确保 Run 只被执行了一次
	once sync.Once
	// 用于确保队列中的任务被正确按顺序执行
	processor *SingleTask.Processor
}
type coinTask struct {
	// 交易对 BTCUSDT
	pair string
}

func CreateBuyCoin(getWave GetWave) BuyCoin {
	bc := BuyCoin{
		coinList:     []string{},
		currentState: buyCoinState{},
		//server:       server,
		//sqlite:       sqlite,
		getWave:   getWave,
		processor: nil,
	}
	bc.processor = SingleTask.NewProcessor(func(i interface{}) {
		task := i.(coinTask)
		bc.runCoinTask(task)
	}, func(i interface{}) {
		ct := i.(coinTask)
		fmt.Printf("task [%s] drop\r\n", ct.pair)
	}, 100)
	go bc.processor.Run()
	bc.updateBuyCoinState()
	return bc
}
func (buyCoin *BuyCoin) UpdateRate(newStopRate float64, newProfitRate float64) {
	buyCoin.getWave.logger.Info(fmt.Sprintf("更新止盈止损比例，新比例止盈 = [%f],止损 = [%f]", newProfitRate, newStopRate))
	if newStopRate != 0 {
		stopRate = newStopRate
	}
	if newProfitRate != 0 {
		profitRate = newProfitRate
	}
}

var updateBuyCoinStateLock = sync.Mutex{}

// 更新状态信息
// - 初始化需要调用
// - 判断是否购买新的交易对也需要判断
// - 更新交易对最后交易情况时（例如买入卖出被触发后）
// test BOME
func (buyCoin *BuyCoin) updateBuyCoinState() error {
	// 这个如果
	updateBuyCoinStateLock.Lock()
	defer updateBuyCoinStateLock.Unlock()
	// 获取 usdt 余额
	fmt.Printf("[updateBuyCoinState] get account info")
	api := UMFAccountAccount.CreateAccountApi(
		UMFAccountAccount.AccountParam{
			RecvWindow: 10000,
		},
	)
	ret, err := buyCoin.getWave.server.Request(api)
	if err != nil {
		buyCoin.getWave.logger.Err(fmt.Sprintf("获取账户信息失败，【%s】", err.Error()))
		return err
	}
	// 请求所有的资产，找到其中 USDT 的资产
	// 这里请求超时了
	account := UMFAccountAccount.ParseResponseToBalance(ret)
	buyCoin.currentState.currentUsdt, _ = strconv.ParseFloat(
		account.GetAssertByName(balanceCoinName).WalletBalance,
		10,
	)
	// 拉取当前挂单情况
	buyCoin.updateTradingInfo(account.Positions)
	// 清洗挂单内容（丢弃并撤销无效订单）
	return buyCoin.updateOpenInfo()
}

// todo test
func (buyCoin *BuyCoin) GetBuyCoinState() (buyCoinState, error) {
	err := buyCoin.updateBuyCoinState()
	return buyCoin.currentState, err
}
func (buyCoin *BuyCoin) SetBuyCoinState() {
	buyCoin.currentState.currentUsdt = 10
	buyCoin.currentState.tradingInfo = []tradingInfo{
		{
			coinName:      "A",
			pair:          "AUSDT",
			createTime:    "12313",
			endTime:       "12313",
			long:          true,
			short:         true,
			longProfit:    0,
			shortProfit:   0,
			longSz:        12,
			shortSz:       12,
			longTradeId:   123,
			shortTradeId:  123,
			longProfitId:  123,
			longStopId:    123,
			shortProfitId: 123,
			shortStopId:   123,
		},
	}
}

// 虽然这里是复数，但是我们只假定这里只有一种币的交易
func (buyCoin *BuyCoin) updateTradingInfo(position []UMFAccountAccount.Position) {
	buyCoin.currentState.tradingInfo = []tradingInfo{}
	//buyCoin.currentState.currentUsdt = 0
	for _, pos := range position {
		symbol := pos.Symbol
		if symbol == balanceCoinName {
		} else {
			// symbol = ETHUSDT
			// symbolName = ETH
			symbolName := symbol[:len(symbol)-len(balanceCoinName)]
			if len(buyCoin.currentState.tradingInfo) == 0 {
				buyCoin.currentState.tradingInfo = []tradingInfo{
					{
						// coinName BTC		pair BTCUSDT
						coinName:      symbolName,
						pair:          symbol,
						createTime:    "",
						endTime:       "",
						long:          false,
						short:         false,
						longProfit:    0,
						shortProfit:   0,
						longSz:        0,
						shortSz:       0,
						longTradeId:   0,
						shortTradeId:  0,
						longProfitId:  0,
						longStopId:    0,
						shortProfitId: 0,
						shortStopId:   0,
					},
				}
			}
			if pos.PositionSide == string(Trade.PositionSideLong) {
				buyCoin.currentState.tradingInfo[0].long = true
				buyCoin.currentState.tradingInfo[0].longProfit, _ = strconv.ParseFloat(pos.UnrealizedProfit, 10)
				buyCoin.currentState.tradingInfo[0].longSz, _ = strconv.ParseFloat(pos.PositionAmt, 10)
			} else if pos.PositionSide == string(Trade.PositionSideShort) {
				buyCoin.currentState.tradingInfo[0].short = true
				buyCoin.currentState.tradingInfo[0].shortProfit, _ = strconv.ParseFloat(pos.UnrealizedProfit, 10)
				buyCoin.currentState.tradingInfo[0].shortSz, _ = strconv.ParseFloat(pos.PositionAmt, 10)
			}
		}
	}
}

type questionOrder struct {
	longOrders  []int64
	shortOrders []int64
	symbol      string
}

func (buyCoin *BuyCoin) updateOpenInfo() error {
	// 这里应该获取所有的 挂单信息
	l := len(buyCoin.currentState.tradingInfo)
	// {
	//		ETHUSDT: 0,
	//		ETHUSDT: 0,
	// }
	coinIndex := make(map[string]int, l)
	for i := 0; i < l; i++ {
		coinIndex[buyCoin.currentState.tradingInfo[i].pair] = i
	}
	api := UMFTradeOpen.CreateOpenApi(
		UMFTradeOpen.OpenParam{
			//Symbol: symbol,
		},
	)
	// 存在问题的 挂单 id，应该需要被撤销
	questionIds := []questionOrder{}
	questiongMap := map[string]int{}
	ret, err := buyCoin.getWave.server.Request(api)
	if err != nil {
		buyCoin.getWave.logger.Info(fmt.Sprintf("获取挂单信息错误 【%s】", err.Error()))
		return err
	}
	order := UMFTradeOpen.ParseResponseToBalance(ret)
	var idx int
	var ok bool
	for _, od := range order {
		//od = {
		//	"orderId": 8389765753007840000,
		//		"symbol": "ETHUSDT",
		//		"status": "NEW",
		//		"clientOrderId": "0csXObtjJh3hBV0dD8ALPk",
		//		"price": "0",
		//		"avgPrice": "0",
		//		"origQty": "1",
		//		"executedQty": "0",
		//		"cumQuote": "0.00000",
		//		"timeInForce": "GTC",
		//		"type": "STOP_MARKET",
		//		"reduceOnly": true,
		//		"closePosition": false,
		//		"side": "SELL",
		//		"positionSide": "LONG",
		//		"stopPrice": "1000",
		//		"workingType": "CONTRACT_PRICE",
		//		"priceProtect": false,
		//		"origType": "STOP_MARKET",
		//		"priceMatch": "NONE",
		//		"selfTradePreventionMode": "NONE",
		//		"goodTillDate": 0,
		//		"time": 1729939238195,
		//		"updateTime": 1729939238195
		//}
		positionSide := Trade.PositionSideType(od.PositionSide)
		if idx, ok = coinIndex[od.Symbol]; ok {
			// 交易对在当前交易信息中
			origType := Trade.OrigType(od.OrigType)
			if Trade.PositionSideLong == positionSide {
				// 开多
				if origType == Trade.OrigTypeStopMarket {
					// 止损
					buyCoin.currentState.tradingInfo[idx].longStopId = od.OrderId
				} else {
					// 止盈
					buyCoin.currentState.tradingInfo[idx].longProfitId = od.OrderId
				}
			} else {
				// 开空
				if origType == Trade.OrigTypeStopMarket {
					// 止损
					buyCoin.currentState.tradingInfo[idx].shortStopId = od.OrderId
				} else {
					// 止盈
					buyCoin.currentState.tradingInfo[idx].shortProfitId = od.OrderId
				}
			}
		} else {
			if _, ok := questiongMap[od.Symbol]; !ok {
				questiongMap[od.Symbol] = len(questionIds)
				questionIds = append(questionIds, questionOrder{
					longOrders:  []int64{},
					shortOrders: []int64{},
					symbol:      od.Symbol,
				})
			}
			if positionSide == Trade.PositionSideLong {
				questionIds[questiongMap[od.Symbol]].longOrders = append(questionIds[questiongMap[od.Symbol]].longOrders, od.OrderId)
			} else {
				questionIds[questiongMap[od.Symbol]].shortOrders = append(questionIds[questiongMap[od.Symbol]].shortOrders, od.OrderId)
			}
		}
	}
	// 这里需要关闭没必要的单，例如现在只剩下开多，但是有 空单 的 止盈止损
	//info := buyCoin.currentState.tradingInfo
	// 放弃下面逻辑，修改为如果多单空单只有一个表示挂单应该清掉了
	for idx, qid := range questionIds {
		if len(qid.longOrders) == 2 {
			questionIds[idx].longOrders = []int64{}
		}
		if len(qid.shortOrders) == 2 {
			questionIds[idx].shortOrders = []int64{}
		}
	}

	// 没有多单，但是有多单的 止盈止损
	//for idx = 0; idx < l; idx++ {
	//	symbol := info[idx].pair
	//	if !info[idx].long {
	//		if info[idx].longProfitId != 0 {
	//			// 问题单
	//			questionIds = append(questionIds, questionOrder{
	//				orderId: info[idx].longProfitId,
	//				symbol:  symbol,
	//			})
	//			buyCoin.currentState.tradingInfo[idx].longProfitId = 0
	//		}
	//		if info[idx].longStopId != 0 {
	//			questionIds = append(questionIds, questionOrder{
	//				orderId: info[idx].longStopId,
	//				symbol:  symbol,
	//			})
	//			buyCoin.currentState.tradingInfo[idx].longStopId = 0
	//		}
	//	} else if !info[idx].short {
	//		if info[idx].shortProfitId != 0 {
	//			// 问题单
	//			questionIds = append(questionIds, questionOrder{
	//				orderId: info[idx].shortProfitId,
	//				symbol:  symbol,
	//			})
	//			buyCoin.currentState.tradingInfo[idx].shortProfitId = 0
	//		}
	//		if info[idx].shortStopId != 0 {
	//			questionIds = append(questionIds, questionOrder{
	//				orderId: info[idx].shortStopId,
	//				symbol:  symbol,
	//			})
	//			buyCoin.currentState.tradingInfo[idx].shortStopId = 0
	//		}
	//	}
	//}
	// 批量 撤销 问题挂单
	//for _, qid := range questionIds {
	//	fmt.Println(qid.orderId, qid.symbol)
	//}
	for _, qid := range questionIds {
		ids := []int64{}
		ids = append(ids, qid.shortOrders...)
		ids = append(ids, qid.longOrders...)
		for _, id := range ids {
			api := UMFTradeOrderDelete.CreateOrderDeleteApi(
				UMFTradeOrderDelete.OrderDeleteParam{
					Symbol:  qid.symbol,
					OrderId: int64(id),
					//OrigClientOrderId: "",
				},
			)
			ret, _ := buyCoin.getWave.server.RequestL(api, true)
			fmt.Printf("[reverse order log] pair = {%s}; orderId = {%d}\r\n", qid.symbol, id)
			fmt.Printf("[reverse order log] %s\r\n", ret)
		}
	}
	return nil
}

// Run 这里新开一个线程
func (buyCoin *BuyCoin) Run() {
	buyCoin.once.Do(func() {
		// 新开线程，执行耗时事务
		go buyCoin.processor.Run()
	})
}

// AddCoin 添加需要判断是否值得购买的币种
// coinName = BTC
func (buyCoin *BuyCoin) AddCoin(coinName string) {
	buyCoin.AddPair(coinName + "USDT")
}

// AddPair 添加需要判断是否值得购买的币种
// pair = BTCUSDT
func (buyCoin *BuyCoin) AddPair(pair string) {
	buyCoin.processor.AddCoin(coinTask{
		pair: pair,
	})
}

// 执行交易对 交易 任务
func (buyCoin *BuyCoin) runCoinTask(task coinTask) {
	// 任务分为几个部分进行
	// ① 收集账户信息
	// 可以拿到当前的
	// 		余额
	//		是否已经开了单
	//		清理问题 挂单
	buyCoin.updateBuyCoinState()
	// 这里后期如果存在多个交易对同时执行，需要循环判断
	if len(buyCoin.currentState.tradingInfo) > 0 && task.pair == buyCoin.currentState.tradingInfo[0].pair {
		// 已经存在当前的交易对 交易 了，不需要再执行了
		buyCoin.getWave.logger.Warn("当前已经持有交易对，停止交易")
		return
	}

	// 获取 当前要进行的交易的 最大开仓数据
	// 设置 当前要进行的交易的 交易杠杆和全仓模式
	leverage := buyCoin.getWave.SetCrossAndMaxLeverage(task.pair)
	if buyCoin.currentState.currentUsdt*float64(leverage) < baseBalance*1.2 {
		// 没钱无法执行
		// todo 这里应该判断是否需要清理掉没用的交易信息
		return
	} else {
		// 获取最后的价格
		price := buyCoin.getLastKLineData(task.pair, 0)
		// 获取 当前要进行的交易的 最小交易量
		detail := buyCoin.calcDearDetail(task.pair, price, leverage)
		fmt.Println(detail)
		buyCoin.getWave.logger.Info(detail.print())
		if detail.ok {
			buyCoin.getWave.logger.Info("================ 交易开始 ================")
			buyCoin.getWave.logger.Info(fmt.Sprintf("可以对[%s]交易", detail.pair))
			buyCoin.getWave.logger.Info(fmt.Sprintf("成交价格大概是 [%f]", price))
			buyCoin.getWave.logger.Info(fmt.Sprintf("开张 [%f]", detail.qty))
			buyCoin.getWave.logger.Info(fmt.Sprintf("开多止盈 [%f] [%f]；止损 [%f] [%f]", detail.longProfit, detail.longProfit/price, detail.longStop, detail.longStop/price))
			buyCoin.getWave.logger.Info(fmt.Sprintf("开空止盈 [%f] [%f]；止损 [%f] [%f]", detail.shortProfit, detail.shortProfit/price, detail.shortStop, detail.shortStop/price))
			detail.toOrder(buyCoin.getWave)
			buyCoin.getWave.logger.Info("================ 交易完成 ================")

			getWaveOrderRecord := Table.GetWaveOrder{
				OrderId:       "",
				Pair:          detail.pair,
				Qty:           detail.qty,
				LongProfitId:  0,
				ShortProfitId: 0,
				LongStopId:    0,
				ShortStopId:   0,
				LongProfit:    detail.longProfit,
				ShortProfit:   detail.shortProfit,
				LongStop:      detail.longStop,
				ShortStop:     detail.shortStop,
				StartTime:     0,
			}
			buyCoin.getWave.sqlite.GetDb().Create(&getWaveOrderRecord)
		}
	}
}

func (buyCoin *BuyCoin) getLastKLineData(pair string, tryTime int) float64 {
	// 最后需要获取一次价格，不然可能因为任务太久远导致价格偏差太大
	api := UMFMarketKLine.CreateKLineApi(UMFMarketKLine.KLineListApiParam{
		Symbol:    pair,
		Interval:  BinaApis.Interval1m, // 使用 1m 即可
		StartTime: 0,
		EndTime:   0,
		Limit:     1, // 获取 1 个 k 柱 即可
	})
	kRet, _ := buyCoin.getWave.server.Request(api)
	kResp := UMFMarketKLine.ParseKLineResponse(kRet)
	if len(kResp) > 0 {
		// 这里无法只获取一个k，最小是10，所以这里应该取最后一个
		k := kResp[0]
		for _, kk := range kResp {
			if kk.OpenTime > k.OpenTime {
				k = kk
			}
		}
		return k.Close
	} else {
		if tryTime > 10 {
			panic(errors.New(fmt.Sprintf("[%s] Get Last KLine Data Fail.", pair)))
		}
		return buyCoin.getLastKLineData(pair, tryTime+1)
	}
}

// 计算交易细节
func (buyCoin BuyCoin) calcDearDetail(pair string, price float64, leverage int) dearDetail {
	detail := dearDetail{
		pair:        pair,
		qty:         0,
		longProfit:  0,
		longStop:    0,
		shortProfit: 0,
		shortStop:   0,
		ok:          false,
	}
	// 获取 当前要进行的交易的 交易细节
	exInfo := buyCoin.getWave.GetExchangeInfoFromDb(pair)
	// 规范化价格
	// 例如 BTC 价格是 60000.0012， FPriceStep = 0.001
	// 格式化后是 60000.002
	price = CoreUtils.Align(price, exInfo.FPriceStep)
	price += exInfo.FPriceStep

	// 当前的 usdt 余额的一半能用来开多和开空，能开多少
	// 当前的钱要乘以开的倍率
	totalUsdt := buyCoin.currentState.currentUsdt * float64(leverage)
	// 保险起见，这里不是除以 2 ，而是除以 2.1，最后取整
	usdt := CoreUtils.Fix(totalUsdt/2.3, 0)
	// 最小需要的交易数量
	qtyLen := CoreUtils.GetPointLen(exInfo.FQtyStep)
	qty := CoreUtils.Fix(usdt/price, qtyLen)
	// 保证四舍五入后的交易量不能低于最低标准
	qty += exInfo.FQtyStep
	if qty < exInfo.FMinQty {
		qty = exInfo.FMinQty
	}
	if qty*price*2 > baseBalance {
		// 可能钱不太够
	}
	// 量太大（感觉不可能触发）
	if qty > exInfo.FMaxQty {
		qty = exInfo.FMaxQty
	}
	// 计算一下当前的交易量需要多少的多少 U，看看余额够不够
	// * 2 是因为同时开了两个方向
	buyCoin.getWave.logger.Info(fmt.Sprintf("current price is %f", price))
	requireUsdt := qty * 2 * price
	// 比最小要求的下单金额大，比当前全部资产高
	if requireUsdt > exInfo.FMinNotional*2 && requireUsdt < totalUsdt {
		fixLen := CoreUtils.GetPointLen(exInfo.FPriceStep)
		detail.qty = qty
		detail.longStop = CoreUtils.Fix(price*(1-stopRate), fixLen)
		detail.longProfit = CoreUtils.Fix(price*(1+profitRate), fixLen)
		detail.shortStop = CoreUtils.Fix(price*(1+stopRate), fixLen)
		detail.shortProfit = CoreUtils.Fix(price*(1-profitRate), fixLen)
		detail.ok = true
	}
	return detail
}

type dearDetail struct {
	// 交易对
	pair string
	// 成交量
	qty float64
	// 开多止盈价格
	longProfit float64
	// 开多止损价格
	longStop float64
	// 开空止盈价格
	shortProfit float64
	// 开空止损价格
	shortStop float64
	// 是否可以下单
	ok bool
}

func (detail dearDetail) print() []string {
	strs := []string{}
	strs = append(strs, "================ dear detail info ================")
	strs = append(strs, fmt.Sprintf("pair: %s", detail.pair))
	strs = append(strs, fmt.Sprintf("qty: %f", detail.qty))
	strs = append(strs, fmt.Sprintf("longStop: %f", detail.longStop))
	strs = append(strs, fmt.Sprintf("longProfit: %f", detail.longProfit))
	strs = append(strs, fmt.Sprintf("shortProfit: %f", detail.shortProfit))
	strs = append(strs, fmt.Sprintf("shortStop: %f", detail.shortStop))
	strs = append(strs, "==================================================")
	for _, s := range strs {
		fmt.Println(s)
	}
	return strs
}

// 去下单和挂单
// 下单的时候去数据库进行登记，方便后期对数据进行收集
func (detail dearDetail) toOrder(wave GetWave) {
	fmt.Println("[to order]")
	// 止盈止损单
	orderParams := []struct {
		price float64
		qty   float64
		order UMFTradeOrder.StopOrder
		desc  string
	}{
		{
			desc:  "open long profit",
			price: detail.longProfit,
			qty:   detail.qty,
			order: UMFTradeOrder.StopOrderLongProfit,
		},
		{
			desc:  "open long stop",
			price: detail.longStop,
			qty:   detail.qty,
			order: UMFTradeOrder.StopOrderLongLoss,
		},
		{
			desc:  "open short stop",
			price: detail.shortStop,
			qty:   detail.qty,
			order: UMFTradeOrder.StopOrderShortLoss,
		},
		{
			desc:  "open short profit",
			price: detail.shortProfit,
			qty:   detail.qty,
			order: UMFTradeOrder.StopOrderShortProfit,
		},
		{
			desc:  "open long",
			qty:   detail.qty,
			order: UMFTradeOrder.StartOrderLongMarket,
		},
		{
			desc:  "open short",
			qty:   detail.qty,
			order: UMFTradeOrder.StartOrderShortMarket,
		},
	}
	// 市价下单
	for _, p := range orderParams {
		fmt.Println(fmt.Sprintf("================= [to order]: %s =================", p.desc))
		param := UMFTradeOrder.OrderParam{
			Symbol: detail.pair,
		}
		if p.price != 0 {
			param.OrderParamStopProfit(p.price, p.qty, p.order)
		} else {
			param.OrderParamOrder(p.qty, p.order, p.price)
		}
		// todo 调试，不发请求
		ret, err := wave.server.Request(UMFTradeOrder.CreateOrderApi(param))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ret)
		fmt.Println("===================================================================")
	}
}
