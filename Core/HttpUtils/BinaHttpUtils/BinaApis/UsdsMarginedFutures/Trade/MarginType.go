package Trade

import (
	"UPureApi/Core/HttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
)

type MarginTypeType string

const (
	MarginTypeISOLATED MarginTypeType = "ISOLATED"
	MarginTypeCROSSED  MarginTypeType = "CROSSED"
)

type MarginTypeParam struct {
	Symbol     string
	MarginType MarginTypeType
}

func (marginTypeParam MarginTypeParam) ToMap() BinaHttpUtils.ParamMap {
	m := BinaHttpUtils.ParamMap{
		"symbol":     BinaApis.CheckEmptyString(marginTypeParam.Symbol),
		"marginType": BinaApis.CheckEmptyString(string(marginTypeParam.MarginType)),
	}
	return m
}

func CreateMarginTypeParam(param MarginTypeParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/marginType",
		HttpMethod:  HttpUtils.HttpPost,
		QueryParams: param,
		BodyParams:  nil,
		Sign:        true,
		Header:      nil,
	}
}
