package UMFMarketPremiumIndex

import (
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
)

// PremiumIndexParam https://www.binance.com/fapi/v1/premiumIndex
// 最新标记价格和资金费率
type PremiumIndexParam struct {
	Symbol string
}

func (premiumIndex PremiumIndexParam) ToMap() BinaHttpUtils.ParamMap {
	m := CoreUtils.ApiParamMap{}
	m.SetNotEmptyString("symbol", premiumIndex.Symbol)
	return BinaHttpUtils.ParamMap(m)
}

func CreatePremiumIndexApi(premiumIndex PremiumIndexParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/premiumIndex",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: premiumIndex,
		BodyParams:  nil,
		Sign:        false,
		Header:      nil,
	}
}
