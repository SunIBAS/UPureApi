package Core

import (
	"fmt"
	"strconv"
)

type ApiParamMap map[string]string

func (apiParamMap ApiParamMap) SetNotEmptyString(key, value string) {
	if value != "" {
		apiParamMap[key] = value
	}
}
func (apiParamMap ApiParamMap) SetNotZeroInt64String(key string, value int64) {
	if value != 0 {
		apiParamMap[key] = strconv.FormatInt(value, 10)
	}
}
func (apiParamMap ApiParamMap) SetNotZeroFloat64String(key string, value float64) {
	if value != 0 {
		apiParamMap[key] = fmt.Sprintf("%f", value)
	}
}

func (apiParamMap ApiParamMap) SetNotZeroDecimal(key string, value float64) {
	if value != 0 {
		//apiParamMap[key] = fmt.Sprintf("%0.0f", value)
		apiParamMap.SetNotZeroFloat64String(key, value)
	}
}
