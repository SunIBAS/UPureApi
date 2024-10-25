package UMFMarketBrackets

type BracketsResponseDataBracket struct {
	Symbol        string `json:"symbol"`
	UpdateTime    int64  `json:"updateTime"`
	NotionalLimit int    `json:"notionalLimit"`
	RiskBrackets  []struct {
		BracketSeq                   int     `json:"bracketSeq"`
		BracketNotionalFloor         int     `json:"bracketNotionalFloor"`
		BracketNotionalCap           int     `json:"bracketNotionalCap"`
		BracketMaintenanceMarginRate float64 `json:"bracketMaintenanceMarginRate"`
		CumFastMaintenanceAmount     float64 `json:"cumFastMaintenanceAmount"`
		MinOpenPosLeverage           int     `json:"minOpenPosLeverage"`
		MaxOpenPosLeverage           int     `json:"maxOpenPosLeverage"`
	} `json:"riskBrackets"`
}
type BracketsResponse struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          struct {
		Brackets []BracketsResponseDataBracket `json:"brackets"`
		Version  string                        `json:"version"`
	} `json:"data"`
	Success bool `json:"success"`
}
