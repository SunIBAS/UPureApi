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
	// 这个是价格精度，例如为 3，表示价格的小数点可以去到 0.001
	PricePrecision        int     `json:"pricePrecision"`
	QuantityPrecision     int     `json:"quantityPrecision"`
	QuoteAsset            string  `json:"quoteAsset"`
	QuotePrecision        int     `json:"quotePrecision"`
	RequiredMarginPercent float64 `json:"requiredMarginPercent"`
	Status                string  `json:"status"`
	TriggerProtect        float64 `json:"triggerProtect"`
	UnderlyingType        string  `json:"underlyingType"`
	// filter 中的选项

	// MinNotional 是最小投入的 U 的数量，无论开的是几倍，都需要使用这个数值对价格求商，向上取整，
	// 例如 RequiredMarginPercent = 5
	// 币 价格为 0.21
	// 5 / 0.21 = 23.8095 -> 向上取整为 24
	// {
	//     "notional": "5",
	//     "filterType": "MIN_NOTIONAL"
	// }
	FMinNotional float64 `json:"minNotional"`

	// 理论上下面的判断是不会用到的

	// minPrice 定义了 price/stopPrice 允许的最小值
	// maxPrice 定义了 price/stopPrice 允许的最大值。
	// tickSize 定义了 price/stopPrice 的步进间隔，即price必须等于minPrice+(tickSize的整数倍) 以上每一项均可为0，为0时代表这一项不再做限制。
	// price >= minPrice
	// price <= maxPrice
	// (price-minPrice) % tickSize == 0
	// 这里的 price 表示下单和止盈止损的 价格
	// {
	//     "filterType": "PRICE_FILTER",
	//     "maxPrice": "100000",
	//     "tickSize": "0.001",
	//     "minPrice": "0.111"
	// }
	FPriceStep float64 // tickSize
	FMaxPrice  float64
	FMinPrice  float64
	// minQty 表示 quantity 允许的最小值.
	// maxQty 表示 quantity 允许的最大值
	// stepSize 表示 quantity允许的步进值。
	// quantity >= minQty
	// quantity <= maxQty
	// (quantity-minQty) % stepSize == 0
	// {
	//     "minQty": "0.1",
	//     "maxQty": "1000000",
	//     "filterType": "LOT_SIZE",
	//     "stepSize": "0.1"
	// }
	FMinQty  float64
	FMaxQty  float64
	FQtyStep float64
}
