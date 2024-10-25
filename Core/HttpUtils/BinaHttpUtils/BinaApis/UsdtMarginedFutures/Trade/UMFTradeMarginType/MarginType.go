package UMFTradeMarginType

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
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

func CreateMarginTypeApi(param MarginTypeParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/marginType",
		HttpMethod:  HttpUtilsCore.HttpPost,
		QueryParams: param,
		BodyParams:  nil,
		Sign:        true,
		Header:      nil,
	}
}
