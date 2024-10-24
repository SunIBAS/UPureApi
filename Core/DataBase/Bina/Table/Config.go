package Table

type Config struct {
	Id             string `json:"symbol" gorm:"primaryKey"`
	Name           string `json:"name"`
	ContentString  string `json:"content"`
	ContentInt64   int64  `json:"contentInt64"`
	ContentFloat64 int64  `json:"contentFloat64"`
}

var (
	// ExchangeInfoSymbolLastUpdate : ExchangeInfoSymbol 数据最后一次更新时间
	ExchangeInfoSymbolLastUpdate Config = Config{
		Id:             "ExchangeInfoSymbolLastUpdate",
		Name:           "ExchangeInfoSymbolLastUpdate",
		ContentString:  "",
		ContentInt64:   0,
		ContentFloat64: 0,
	}
	// ExchangeInfoSymbolLiveTime : ExchangeInfoSymbol 存活时间(十天)
	ExchangeInfoSymbolLiveTime int64 = 1000 * 60 * 60 * 24 * 10
	KLineDataSymbols                 = Config{
		Id:             "",
		Name:           "symbol",
		ContentString:  "",
		ContentInt64:   0,
		ContentFloat64: 0,
	}
	BracketsConfigSyncComplete = "true"
	BracketsConfig             = Config{
		Id:             "brackets-sync",
		Name:           "brackets",
		ContentString:  "",
		ContentInt64:   0,
		ContentFloat64: 0,
	}
)
