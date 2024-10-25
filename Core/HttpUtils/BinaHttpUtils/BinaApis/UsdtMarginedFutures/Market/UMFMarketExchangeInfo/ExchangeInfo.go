package UMFMarketExchangeInfo

import (
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/DataBase/Bina/Table"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"encoding/json"
	"fmt"
)

type ExchangeInfoParam struct {
}

func (exchangeInfo ExchangeInfoParam) ToMap() BinaHttpUtils.ParamMap {
	return map[string]string{
		//"showall": "true",
	}
}

// CreateExchangeInfoApi https://binance-docs.github.io/apidocs/futures/cn/#api
func CreateExchangeInfoApi(exchangeInfo ExchangeInfoParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/exchangeInfo",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: exchangeInfo,
		BodyParams:  nil,
		Sign:        false,
		Header:      nil,
		NoTimeStamp: true,
	}
}

func ParseExchangeInfoResponse(jsonData string) ExchangeInfoResponse {
	// 解析 JSON 数据
	var exchangeInfo ExchangeInfoResponse
	err := json.Unmarshal([]byte(jsonData), &exchangeInfo)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return ExchangeInfoResponse{}
	}

	return exchangeInfo
}

func ParseExchangeInfoResponseTable(symbol Symbol) Table.ExchangeInfoSymbol {
	return Table.ExchangeInfoSymbol{
		BaseAsset:             symbol.BaseAsset,
		BaseAssetPrecision:    symbol.BaseAssetPrecision,
		ContractType:          symbol.ContractType,
		DeliveryDate:          symbol.DeliveryDate,
		LiquidationFee:        CoreUtils.StringToFloat64(symbol.LiquidationFee),
		MaintMarginPercent:    CoreUtils.StringToFloat64(symbol.MaintMarginPercent),
		MarginAsset:           symbol.MarginAsset,
		MarketTakeBound:       CoreUtils.StringToFloat64(symbol.MarketTakeBound),
		MaxMoveOrderLimit:     symbol.MaxMoveOrderLimit,
		OnboardDate:           symbol.OnboardDate,
		Pair:                  symbol.Pair,
		PricePrecision:        symbol.PricePrecision,
		QuantityPrecision:     symbol.QuantityPrecision,
		QuoteAsset:            symbol.QuoteAsset,
		QuotePrecision:        symbol.QuotePrecision,
		RequiredMarginPercent: CoreUtils.StringToFloat64(symbol.RequiredMarginPercent),
		Status:                symbol.Status,
		Symbol:                symbol.Symbol,
		TriggerProtect:        CoreUtils.StringToFloat64(symbol.TriggerProtect),
		UnderlyingType:        symbol.UnderlyingType,
	}
}
