package Bina

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketBrackets"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"fmt"
	"testing"
)

var url = "https://www.binance.com/bapi/futures/v1/friendly/future/common/brackets"

func TestBarkets(t *testing.T) {
	proxyOption := HttpUtilsCore.SetProxy("http", "localhost", "7890")
	request := HttpUtilsCore.CreateRequest(proxyOption)
	ret, _ := request.ToRequest(HttpUtilsCore.HttpGet, url, HttpUtilsCore.DefaultHeader, nil)
	brackets := UMFMarketBrackets.PParseBracketResponseTable(ret)
	fmt.Println(brackets)
}
