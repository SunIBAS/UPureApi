package Trade

import (
	"UPureApi/Core"
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
)

type LeverAgeParam struct {
	Symbol   string
	Leverage int64
}

func (leverAgeParam LeverAgeParam) ToMap() BinaHttpUtils.ParamMap {
	m := Core.ApiParamMap{
		"symbol": BinaApis.CheckEmptyString(leverAgeParam.Symbol),
	}
	m.SetNotZeroInt64String("leverage", leverageLimit.Get(leverAgeParam.Leverage))
	return BinaHttpUtils.ParamMap(m)
}

var leverageLimit = BinaApis.IntRangeParam{
	Max:     125,
	Min:     1,
	Default: 2,
}

func CreateLeverAgeParam(param LeverAgeParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/leverage",
		HttpMethod:  HttpUtils.HttpPost,
		QueryParams: param,
		BodyParams:  nil,
		Sign:        true,
		Header:      nil,
	}
}
