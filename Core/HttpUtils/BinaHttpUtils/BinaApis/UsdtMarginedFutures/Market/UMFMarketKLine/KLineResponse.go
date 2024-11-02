package UMFMarketKLine

import (
	"UPureApi/Core/CoreUtils"
	"encoding/json"
	"fmt"
)

type KLineResponse struct {
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

func ParseKLineResponse(jsonData string) []KLineResponse {
	var rawItems [][]interface{}
	err := json.Unmarshal([]byte(jsonData), &rawItems)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}
	// 准备一个用于存放结构化数据的切片
	var items []KLineResponse
	for _, rawItem := range rawItems {
		// 构造 Data 结构体，转换各个字段
		item := KLineResponse{
			OpenTime:          int64(rawItem[0].(float64)),
			Open:              CoreUtils.StringToFloat64(rawItem[1].(string)),
			High:              CoreUtils.StringToFloat64(rawItem[2].(string)),
			Low:               CoreUtils.StringToFloat64(rawItem[3].(string)),
			Close:             CoreUtils.StringToFloat64(rawItem[4].(string)),
			Vol:               CoreUtils.StringToFloat64(rawItem[5].(string)),
			CloseTime:         int64(rawItem[6].(float64)),
			Turnover:          CoreUtils.StringToFloat64(rawItem[7].(string)),
			NumberTranslation: rawItem[8].(float64),
			BuyVol:            CoreUtils.StringToFloat64(rawItem[9].(string)),
			BuyTurnover:       CoreUtils.StringToFloat64(rawItem[10].(string)),
			Ignore:            rawItem[11].(string),
		}
		items = append(items, item)
	}
	return items
}
