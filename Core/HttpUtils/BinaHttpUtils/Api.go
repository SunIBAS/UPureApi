package BinaHttpUtils

import (
	"UPureApi/Core/HttpUtils"
	"net/http"
)

type Params interface {
	ToMap() ParamMap
}

type ParamMap map[string]string

func (p ParamMap) ToString() string {
	return HttpUtils.Params2string(p, HttpUtils.ToString)
}

type Api struct {
	NoTimeStamp bool
	Path        string
	HttpMethod  HttpUtils.HttpMethod
	QueryParams Params
	BodyParams  Params
	Sign        bool
	Header      http.Header
}
