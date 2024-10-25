package UMFMarketExchangeInfo

// ExchangeInfoResponse 结构体定义
type ExchangeInfoResponse struct {
	Assets          []Asset       `json:"assets"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	FuturesType     string        `json:"futuresType"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	ServerTime      int64         `json:"serverTime"`
	Symbols         []Symbol      `json:"symbols"`
	Timezone        string        `json:"timezone"`
}

type RateLimit struct {
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
	RateLimitType string `json:"rateLimitType"`
}

type Asset struct {
	Asset             string `json:"asset"`
	MarginAvailable   bool   `json:"marginAvailable"`
	AutoAssetExchange string `json:"autoAssetExchange"`
}

type Symbol struct {
	BaseAsset          string   `json:"baseAsset"`
	BaseAssetPrecision int      `json:"baseAssetPrecision"`
	ContractType       string   `json:"contractType"`
	DeliveryDate       int64    `json:"deliveryDate"`
	Filters            []Filter `json:"filters"`
	LiquidationFee     string   `json:"liquidationFee"`
	MaintMarginPercent string   `json:"maintMarginPercent"`
	MarginAsset        string   `json:"marginAsset"`
	MarketTakeBound    string   `json:"marketTakeBound"`
	MaxMoveOrderLimit  int      `json:"maxMoveOrderLimit"`
	OnboardDate        int64    `json:"onboardDate"`
	OrderTypes         []string `json:"OrderType"`
	Pair               string   `json:"pair"`
	PricePrecision     int      `json:"pricePrecision"`
	QuantityPrecision  int      `json:"quantityPrecision"`
	QuoteAsset         string   `json:"quoteAsset"`
	QuotePrecision     int      `json:"quotePrecision"`
	// 这个是最小投入的 U 的数量，无论开的是几倍，都需要使用这个数值对价格求商，向上取整，
	// 例如 RequiredMarginPercent = 5
	// 币 价格为 0.21
	// 5 / 0.21 = 23.8095 -> 向上取整为 24
	RequiredMarginPercent string   `json:"requiredMarginPercent"`
	Status                string   `json:"status"`
	Symbol                string   `json:"symbol"`
	TimeInForce           []string `json:"timeInForce"`
	TriggerProtect        string   `json:"triggerProtect"`
	UnderlyingSubType     []string `json:"underlyingSubType"`
	UnderlyingType        string   `json:"underlyingType"`
}

type Filter struct {
	FilterType        string `json:"filterType"`
	MaxPrice          string `json:"maxPrice,omitempty"`
	MinPrice          string `json:"minPrice,omitempty"`
	TickSize          string `json:"tickSize,omitempty"`
	MaxQty            string `json:"maxQty,omitempty"`
	MinQty            string `json:"minQty,omitempty"`
	StepSize          string `json:"stepSize,omitempty"`
	Limit             int    `json:"limit,omitempty"`
	Notional          string `json:"notional,omitempty"`
	MultiplierUp      string `json:"multiplierUp,omitempty"`
	MultiplierDown    string `json:"multiplierDown,omitempty"`
	MultiplierDecimal string `json:"multiplierDecimal,omitempty"`
}
