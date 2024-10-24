package Utils

import (
	"encoding/json"
	"fmt"
)

type PremiumIndexKLine struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	InterestRate         string `json:"interestRate"`
	Time                 int64  `json:"time"`
}

// ParsePremiumIndexReturn
//
//	[{
//		"symbol":"REZUSDT",
//		"markPrice":"0.04271000",
//		"indexPrice":"0.04269813",
//		"estimatedSettlePrice":"0.04287215",
//		"lastFundingRate":"0.00005000",
//		"interestRate":"0.00005000",
//		"nextFundingTime":1729526400000,
//		"time":1729517189000
//	}]
func ParsePremiumIndexReturn(str string) []PremiumIndexKLine {
	var pKLines []PremiumIndexKLine
	err := json.Unmarshal([]byte(str), &pKLines)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}
	return pKLines
}
