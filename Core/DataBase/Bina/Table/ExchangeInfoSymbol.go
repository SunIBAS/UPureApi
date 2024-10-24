package Table

type ExchangeInfoSymbol struct {
	Symbol             string  `json:"symbol" gorm:"primaryKey"`
	BaseAsset          string  `json:"baseAsset"`
	BaseAssetPrecision int     `json:"baseAssetPrecision"`
	ContractType       string  `json:"contractType"`
	DeliveryDate       int64   `json:"deliveryDate"`
	LiquidationFee     float64 `json:"liquidationFee"`
	MaintMarginPercent float64 `json:"maintMarginPercent"`
	MarginAsset        string  `json:"marginAsset"`
	MarketTakeBound    float64 `json:"marketTakeBound"`
	MaxMoveOrderLimit  int     `json:"maxMoveOrderLimit"`
	OnboardDate        int64   `json:"onboardDate"`
	Pair               string  `json:"pair"`
	PricePrecision     int     `json:"pricePrecision"`
	QuantityPrecision  int     `json:"quantityPrecision"`
	QuoteAsset         string  `json:"quoteAsset"`
	QuotePrecision     int     `json:"quotePrecision"`
	// 这个是最小投入的 U 的数量，无论开的是几倍，都需要使用这个数值对价格求商，向上取整，
	// 例如 RequiredMarginPercent = 5
	// 币 价格为 0.21
	// 5 / 0.21 = 23.8095 -> 向上取整为 24
	RequiredMarginPercent float64 `json:"requiredMarginPercent"`
	Status                string  `json:"status"`
	TriggerProtect        float64 `json:"triggerProtect"`
	UnderlyingType        string  `json:"underlyingType"`
}
