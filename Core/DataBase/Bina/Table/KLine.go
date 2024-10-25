package Table

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
	"fmt"
)

type KLine struct {
	Id                string  `json:"id" gorm:"primaryKey"`
	OpenTime          int64   `json:"openTime"`
	Open              float64 `json:"open"`
	High              float64 `json:"high"`
	Low               float64 `json:"low"`
	Close             float64 `json:"close"`
	Vol               float64 `json:"vol"` // 成交量
	CloseTime         int64   `json:"closeTime"`
	Turnover          float64 `json:"turnover"`          // 成交金额
	NumberTranslation float64 `json:"numberTranslation"` // 成交笔数
	BuyVol            float64 `json:"BuyVol"`            // 买入量
	BuyTurnover       float64 `json:"buyTurnover"`       // 买入金额
	Ignore            string  `json:"ignore"`            // 随机数
}

func ApiKLine2KLine(symbol string, k UMFMarketKLine.KLineResponse) KLine {
	return KLine{
		Id:                fmt.Sprintf("%s_%d", symbol, k.OpenTime),
		OpenTime:          k.OpenTime,
		Open:              k.Open,
		High:              k.High,
		Low:               k.Low,
		Close:             k.Close,
		Vol:               k.Vol,
		CloseTime:         k.CloseTime,
		Turnover:          k.Turnover,
		NumberTranslation: k.NumberTranslation,
		BuyVol:            k.BuyVol,
		BuyTurnover:       k.BuyTurnover,
		Ignore:            k.Ignore,
	}
}
