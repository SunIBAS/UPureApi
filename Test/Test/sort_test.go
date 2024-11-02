package Test

import (
	"UPureApi/Bina/App/GetWave/Model"
	"fmt"
	"math"
	"sort"
	"testing"
)

type StockVolatility struct {
	Name       string
	Volatility float64
}

// findAndSortVolatileStocks 按波动性对股票进行排序
func findAndSortVolatileStocks(stockData Model.KLineSortInfoArr) []StockVolatility {
	var volatilities []StockVolatility

	for _, sd := range stockData {
		stock := sd.Symbol
		prices := sd.Rate
		// 确保有足够的数据点
		if len(prices) < 2 {
			continue
		}

		// 计算各时间段的波动幅度
		vol15 := math.Abs(1 - prices[0])
		vol30 := math.Abs(1 - prices[1])
		vol45 := math.Abs(1 - prices[2])

		// 检查波动是否持续（短期波动大于长期波动）
		if vol15 > vol30 && vol30 > vol45 {
			// 将符合条件的股票添加到volatilities列表
			volatilities = append(volatilities, StockVolatility{
				Name:       stock,
				Volatility: vol15, // 以15分钟的波动幅度作为排序标准
			})
		}
	}

	// 根据波动性从高到低排序
	sort.Slice(volatilities, func(i, j int) bool {
		return volatilities[i].Volatility > volatilities[j].Volatility
	})

	return volatilities
}

func TestSort(t *testing.T) {
	// 假设的股票数据结构
	data := Model.KLineSortInfoArr{
		{
			Symbol:    "A",
			Rate:      []float64{0.99, 0.94, 0.93},
			Interval:  "15m",
			LastPrice: 0,
		},
		{
			Symbol:    "B",
			Rate:      []float64{0.99, 0.94, 0.93},
			Interval:  "15m",
			LastPrice: 0,
		},
	}

	volatilities := findAndSortVolatileStocks(data)

	// 输出排序后的股票列表
	fmt.Println("Sorted stocks by volatility (highest first):")
	for _, stock := range volatilities {
		fmt.Printf("Stock: %s, Volatility: %.2f%%\n", stock.Name, stock.Volatility*100)
	}

}
