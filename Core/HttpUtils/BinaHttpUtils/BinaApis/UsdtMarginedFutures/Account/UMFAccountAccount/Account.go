package UMFAccountAccount

import (
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"encoding/json"
	"fmt"
)

// Create Time : 19:36

type AccountParam struct {
	RecvWindow int64
}

func (param AccountParam) ToMap() BinaHttpUtils.ParamMap {
	//return map[string]string{}
	p := CoreUtils.ApiParamMap{}
	p.SetNotZeroInt64String("recvWindow", param.RecvWindow)
	return BinaHttpUtils.ParamMap(p)
}

func CreateAccountApi(param AccountParam) BinaHttpUtils.Api {
	return BinaHttpUtils.Api{
		NoTimeStamp: false,
		Path:        "/fapi/v3/account",
		HttpMethod:  HttpUtilsCore.HttpGet,
		QueryParams: param,
		BodyParams:  nil,
		Sign:        true,
		Header:      nil,
	}
}

type Assert struct {
	Asset                  string `json:"asset"`
	WalletBalance          string `json:"walletBalance"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	MarginBalance          string `json:"marginBalance"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	AvailableBalance       string `json:"availableBalance"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
	UpdateTime             int64  `json:"updateTime"`
}
type Position struct {
	Symbol           string `json:"symbol"`
	PositionSide     string `json:"positionSide"`
	PositionAmt      string `json:"positionAmt"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	IsolatedMargin   string `json:"isolatedMargin"`
	Notional         string `json:"notional"`
	IsolatedWallet   string `json:"isolatedWallet"`
	InitialMargin    string `json:"initialMargin"`
	MaintMargin      string `json:"maintMargin"`
	UpdateTime       int    `json:"updateTime"`
}
type AccountResponse struct {
	TotalInitialMargin          string     `json:"totalInitialMargin"`
	TotalMaintMargin            string     `json:"totalMaintMargin"`
	TotalWalletBalance          string     `json:"totalWalletBalance"`
	TotalUnrealizedProfit       string     `json:"totalUnrealizedProfit"`
	TotalMarginBalance          string     `json:"totalMarginBalance"`
	TotalPositionInitialMargin  string     `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin string     `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     string     `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             string     `json:"totalCrossUnPnl"`
	AvailableBalance            string     `json:"availableBalance"`
	MaxWithdrawAmount           string     `json:"maxWithdrawAmount"`
	Assets                      []Assert   `json:"assets"`
	Positions                   []Position `json:"positions"`
}

func ParseResponseToBalance(str string) AccountResponse {
	var ar AccountResponse
	err := json.Unmarshal([]byte(str), &ar)
	if err != nil {
		fmt.Println(err.Error())
	}
	return ar
}

func (accountResponse AccountResponse) GetAssertByName(symbol string) Assert {
	assert := Assert{}
	for _, ast := range accountResponse.Assets {
		if ast.Asset == symbol {
			return ast
		}
	}
	return assert
}
