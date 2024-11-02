package UMFTradeOrderDelete

import (
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
)

// Create Time : 21:21

type OrderDeleteParam struct {
	Symbol            string
	OrderId           int64
	OrigClientOrderId string
	//recvWindow        int64
}

func (param OrderDeleteParam) ToMap() BinaHttpUtils.ParamMap {
	m := CoreUtils.ApiParamMap{
		"symbol": BinaApis.CheckEmptyString(param.Symbol),
	}
	m.SetNotZeroInt64String("orderId", param.OrderId)
	m.SetNotEmptyString("origClientOrderId", param.OrigClientOrderId)
	return BinaHttpUtils.ParamMap(m)
}

func CreateOrderDeleteApi(param OrderDeleteParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		NoTimeStamp: false,
		Path:        "/fapi/v1/order",
		HttpMethod:  HttpUtilsCore.HttpDelete,
		//HttpMethod:  HttpUtilsCore.HttpPost,
		QueryParams: param,
		BodyParams:  nil,
		Sign:        true,
		//		Sign:        false,
		Header: nil,
	}
}

// type OrderDeleteResponse struct {}

// func ParseResponseToBalance(str string) OrderDeleteResponse {
// 	var resp OrderDeleteResponse
// 	err := json.Unmarshal([]byte(str), &resp)
// 	if err!=nil {
// 		fmt.Println(err.Error())
// 	}
// 	return resp
// }
