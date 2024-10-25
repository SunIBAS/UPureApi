package Model

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
)

type KLineSortInfo struct {
	Symbol string
	// 涨幅百分比
	Rate      []float64
	Interval  BinaApis.Interval
	LastPrice float64
}
type KLineSortInfoArr []KLineSortInfo

func BuildKLineSortInfo(symbol string, line []UMFMarketKLine.KLineResponse) KLineSortInfo {
	l := len(line)
	kInfo := KLineSortInfo{
		Symbol:    symbol,
		Rate:      make([]float64, l),
		Interval:  BinaApis.Interval15m,
		LastPrice: 0,
	}
	lastK := line[l-1]
	kInfo.LastPrice = lastK.Close
	for i := l - 1; i > 0; i-- {
		kInfo.Rate[l-i-1] = lastK.Close/line[i].Open - 1
	}
	return kInfo
}
