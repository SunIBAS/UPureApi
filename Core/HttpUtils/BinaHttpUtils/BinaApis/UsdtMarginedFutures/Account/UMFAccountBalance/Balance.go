package UMFAccountBalance

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"encoding/json"
	"fmt"
)

type BalanceParam struct{}

func (BalanceParam BalanceParam) ToMap() BinaHttpUtils.ParamMap {
	return map[string]string{}
}
func CreateBalanceApi(params BalanceParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		NoTimeStamp: false,
		Path:        "/fapi/v3/balance",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: params,
		BodyParams:  nil,
		Sign:        true,
		Header:      nil,
	}
}

type BalanceResponse struct {
	AccountAlias       string `json:"accountAlias"`
	Asset              string `json:"asset"`
	Balance            string `json:"balance"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
	AvailableBalance   string `json:"availableBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
	MarginAvailable    bool   `json:"marginAvailable"`
	UpdateTime         int64  `json:"updateTime"`
}

func ParseResponseToBalance(str string) []BalanceResponse {
	var bResp []BalanceResponse
	err := json.Unmarshal([]byte(str), &bResp)
	if err != nil {
		fmt.Println(err.Error())
	}
	return bResp
}

type BalanceResponseArr []BalanceResponse

func (bArr BalanceResponseArr) GetBuySymbolName(symbol string) BalanceResponse {
	for _, b := range bArr {
		if b.Asset == symbol {
			return b
		}
	}
	return BalanceResponse{
		AccountAlias:       "",
		Asset:              "",
		Balance:            "0",
		CrossWalletBalance: "",
		CrossUnPnl:         "",
		AvailableBalance:   "",
		MaxWithdrawAmount:  "",
		MarginAvailable:    false,
		UpdateTime:         0,
	}
}
