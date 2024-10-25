package UMFMarketKLine

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"strconv"
)

type KLineListApiParam struct {
	Symbol    string
	Interval  BinaApis.Interval
	StartTime int64
	EndTime   int64
	Limit     int64
}

var KLineListLimit = BinaApis.IntRangeParam{
	Max:     1500,
	Min:     10,
	Default: 500,
}

func (kLineApiParam KLineListApiParam) ToMap() BinaHttpUtils.ParamMap {
	m := map[string]string{
		"symbol":   BinaApis.CheckEmptyString(kLineApiParam.Symbol),
		"interval": BinaApis.CheckEmptyString(string(kLineApiParam.Interval)),
		"limit":    strconv.FormatInt(KLineListLimit.Get(kLineApiParam.Limit), 10),
	}
	if kLineApiParam.StartTime > 1000 {
		m["startTime"] = strconv.FormatInt(kLineApiParam.StartTime, 10)
	}
	if kLineApiParam.EndTime > 1000 {
		m["endTime"] = strconv.FormatInt(kLineApiParam.EndTime, 10)
	}
	return m
}

func CreateKLineApi(param KLineListApiParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		Path:        "/fapi/v1/klines",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: param,
		BodyParams:  BinaApis.EmptyParams{},
		Sign:        false,
		Header:      HttpUtilsCore.DefaultHeader,
	}
}
