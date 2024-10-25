package UMFMarketPremiumIndex

import (
	"encoding/json"
	"fmt"
)

type PremiumIndexResponse struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	InterestRate         string `json:"interestRate"`
	Time                 int64  `json:"time"`
}

func ParsePremiumIndexResponse(str string) []PremiumIndexResponse {
	var pKLines []PremiumIndexResponse
	err := json.Unmarshal([]byte(str), &pKLines)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}
	return pKLines
}
