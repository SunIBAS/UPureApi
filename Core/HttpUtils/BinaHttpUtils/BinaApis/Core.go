package BinaApis

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"errors"
)

type IntRangeParam struct {
	Max     int64
	Min     int64
	Default int64
}

func (intRangeParam IntRangeParam) Get(value int64) int64 {
	if value > intRangeParam.Max {
		return intRangeParam.Max
	} else if value < intRangeParam.Min {
		return intRangeParam.Min
	} else {
		return value
	}
}

func CheckEmptyString(s string) string {
	if s == "" {
		panic(errors.New("symbol is empty"))
	}
	return s
}

type Interval string

const (
	Interval1m  Interval = "1m"
	Interval3m  Interval = "3m"
	Interval5m  Interval = "5m"
	Interval15m Interval = "15m"
	Interval30m Interval = "30m"
	Interval1h  Interval = "1h"
	Interval2h  Interval = "2h"
	Interval4h  Interval = "4h"
	Interval6h  Interval = "6h"
	Interval8h  Interval = "8h"
	Interval12h Interval = "12h"
	Interval1d  Interval = "1d"
	Interval3d  Interval = "3d"
	Interval1w  Interval = "1w"
	Interval1M  Interval = "1M"
)

type EmptyParams struct {
}

func (EmptyParams) ToMap() BinaHttpUtils.ParamMap {
	return map[string]string{}
}
func CreateParam() {

}
