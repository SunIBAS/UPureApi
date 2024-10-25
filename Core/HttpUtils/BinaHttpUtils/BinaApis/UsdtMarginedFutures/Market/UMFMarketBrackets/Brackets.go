package UMFMarketBrackets

import (
	"UPureApi/Core/DataBase/Bina/Table"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"encoding/json"
)

/* 获取所有交易对的最大最小交易杠杆 */

// BracketsParam https://www.binance.com/bapi/futures/v1/friendly/future/common/brackets
// https://www.binance.com/zh-CN/futures/trading-rules/perpetual/leverage-margin
type BracketsParam struct{}

func (bracketsParam BracketsParam) ToMap() BinaHttpUtils.ParamMap {
	return map[string]string{}
}

func CreateBracketsApi(bracketsParam BracketsParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		NoTimeStamp: false,
		Path:        "https://www.binance.com/bapi/futures/v1/friendly/future/common/brackets",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: bracketsParam,
		BodyParams:  nil,
		Sign:        false,
		Header:      nil,
	}
}

func ParseBracketsResponse(str string) BracketsResponse {
	var brObj BracketsResponse
	if err := json.Unmarshal([]byte(str), &brObj); err != nil {
		panic(err)
	}
	return brObj
}

func ParseBracketResponseTable(bracket BracketsResponseDataBracket) Table.Brackets {
	return Table.Brackets{
		Symbol:                       bracket.Symbol,
		UpdateTime:                   bracket.UpdateTime,
		NotionalLimit:                bracket.NotionalLimit,
		BracketSeq:                   bracket.RiskBrackets[0].BracketSeq,
		BracketNotionalFloor:         bracket.RiskBrackets[0].BracketNotionalFloor,
		BracketNotionalCap:           bracket.RiskBrackets[0].BracketNotionalCap,
		BracketMaintenanceMarginRate: bracket.RiskBrackets[0].BracketMaintenanceMarginRate,
		CumFastMaintenanceAmount:     bracket.RiskBrackets[0].CumFastMaintenanceAmount,
		MinOpenPosLeverage:           bracket.RiskBrackets[0].MinOpenPosLeverage,
		MaxOpenPosLeverage:           bracket.RiskBrackets[0].MaxOpenPosLeverage,
	}
}

func PParseBracketResponseTable(str string) []Table.Brackets {
	brObj := ParseBracketsResponse(str)
	l := len(brObj.Data.Brackets)
	brackets := make([]Table.Brackets, l)
	for i := 0; i < l; i++ {
		brackets[i] = ParseBracketResponseTable(brObj.Data.Brackets[i])
	}
	return brackets
}
