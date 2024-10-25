package BinaHttpUtils

import (
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"net/http"
)

type Params interface {
	ToMap() ParamMap
}

type ParamMap map[string]string

func (p ParamMap) ToString() string {
	return HttpUtilsCore.Params2string(p, HttpUtilsCore.ToString)
}

type Api struct {
	NoTimeStamp bool
	Path        string
	HttpMethod  HttpUtilsCore.HttpMethod
	QueryParams Params
	BodyParams  Params
	Sign        bool
	Header      http.Header
}
