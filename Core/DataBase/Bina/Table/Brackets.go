package Table

import "encoding/json"

type Brackets struct {
	Symbol        string `json:"symbol"`
	UpdateTime    int64  `json:"updateTime"`
	NotionalLimit int    `json:"notionalLimit"`
	// 只抽取第一层
	BracketSeq                   int     `json:"bracketSeq"`
	BracketNotionalFloor         int     `json:"bracketNotionalFloor"`
	BracketNotionalCap           int     `json:"bracketNotionalCap"`
	BracketMaintenanceMarginRate float64 `json:"bracketMaintenanceMarginRate"`
	CumFastMaintenanceAmount     int     `json:"cumFastMaintenanceAmount"`
	MinOpenPosLeverage           int     `json:"minOpenPosLeverage"`
	MaxOpenPosLeverage           int     `json:"maxOpenPosLeverage"`
}

type BracketsResponseObject struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          struct {
		Brackets []struct {
			Symbol        string `json:"symbol"`
			UpdateTime    int64  `json:"updateTime"`
			NotionalLimit int    `json:"notionalLimit"`
			RiskBrackets  []struct {
				BracketSeq                   int     `json:"bracketSeq"`
				BracketNotionalFloor         int     `json:"bracketNotionalFloor"`
				BracketNotionalCap           int     `json:"bracketNotionalCap"`
				BracketMaintenanceMarginRate float64 `json:"bracketMaintenanceMarginRate"`
				CumFastMaintenanceAmount     int     `json:"cumFastMaintenanceAmount"`
				MinOpenPosLeverage           int     `json:"minOpenPosLeverage"`
				MaxOpenPosLeverage           int     `json:"maxOpenPosLeverage"`
			} `json:"riskBrackets"`
		} `json:"brackets"`
		Version string `json:"version"`
	} `json:"data"`
	Success bool `json:"success"`
}

func ParseBracketsResponseObject(str string) []Brackets {
	var brObj BracketsResponseObject
	if err := json.Unmarshal([]byte(str), &brObj); err != nil {
		panic(err)
	}
	var brackets = make([]Brackets, len(brObj.Data.Brackets))
	for idx, b := range brObj.Data.Brackets {
		brackets[idx] = Brackets{
			Symbol:                       b.Symbol,
			UpdateTime:                   b.UpdateTime,
			NotionalLimit:                b.NotionalLimit,
			BracketSeq:                   b.RiskBrackets[0].BracketSeq,
			BracketNotionalFloor:         b.RiskBrackets[0].BracketNotionalFloor,
			BracketNotionalCap:           b.RiskBrackets[0].BracketNotionalCap,
			BracketMaintenanceMarginRate: b.RiskBrackets[0].BracketMaintenanceMarginRate,
			CumFastMaintenanceAmount:     b.RiskBrackets[0].CumFastMaintenanceAmount,
			MinOpenPosLeverage:           b.RiskBrackets[0].MinOpenPosLeverage,
			MaxOpenPosLeverage:           b.RiskBrackets[0].MaxOpenPosLeverage,
		}
	}
	return brackets
}
