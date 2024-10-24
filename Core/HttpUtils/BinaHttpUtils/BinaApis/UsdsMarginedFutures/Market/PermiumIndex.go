package Market

import (
	"UPureApi/Core"
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
)

// https://www.binance.com/fapi/v1/premiumIndex
type PremiumIndex struct {
	Symbol string
}

func (premiumIndex PremiumIndex) ToMap() BinaHttpUtils.ParamMap {
	m := Core.ApiParamMap{}
	m.SetNotEmptyString("symbol", premiumIndex.Symbol)
	return BinaHttpUtils.ParamMap(m)
}

func CreatePremiumIndexApi(premiumIndex PremiumIndex) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/premiumIndex",
		HttpMethod:  HttpUtils.HttpGet,
		QueryParams: premiumIndex,
		BodyParams:  nil,
		Sign:        false,
		Header:      nil,
	}
}
