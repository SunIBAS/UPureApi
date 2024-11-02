package Utils

import (
	"UPureApi/Bina/App/GetWave/Model"
	"sort"
)

type StockVolatility struct {
	Name string
	Wave float64
}

// FindAndSortVolatileStocks 按波动性对股票进行排序
// checkLen 检查的长度
// okLen 	满足爬坡的长度
func FindAndSortVolatileStocks(stockData Model.KLineSortInfoArr, checkLen int, okLen int) []StockVolatility {
	var volatilities []StockVolatility

	bigTime := 1
	var lastOne float64 = -1
	tmp := make([]float64, checkLen)
	for _, sd := range stockData {
		bigTime = 1
		lastOne = -1
		if sd.Rate == nil || len(sd.Rate) < checkLen {
			continue
		}
		// Rate = [
		//    -0.005377950403346321,
		//    0.002408912978018707,
		//    -0.0044856459330142595,
		//    -0.008636092912447824,
		//    -0.01187295933511423,
		//    -0.012752075919335693,
		//    -0.012166172106825024,
		//    -0.01157957244655583,
		//    -0.007453786523553996,
		//    -0.007749627421758554
		//]
		// tmp = [
		//    1,				   //
		//    -0.4479239854126973, // 变小
		//    0.8340809409888118,  // 变小
		//    1.6058334987757026,  // 变大[√]
		//    2.2077108274792794,  // 变大[√]
		//    2.3711776723342357,  // 变大[√]
		//    2.262232113419058,   // 变小
		//    2.153157165479004,   // 变小
		//    1.3859901941297241,  // 变小
		//    1.4410001655903157   // 变小
		//]
		// ① 计算坡度，即和第一个的比值（这里同时可以消去正负号）
		// ② 从第一个开始计算后续变大的量
		//		例如 1 后面是 -0.44 则变小
		//			一直到 1.60 变大
		//			2.20 变大（这里是和 1.60 对比，不是和 1 对比）
		for i := 0; i < checkLen; i++ {
			tmp[i] = sd.Rate[i] / sd.Rate[0]
			if i != 0 {
				if tmp[i] > lastOne {
					lastOne = tmp[i]
					bigTime++
				}
			} else {
				lastOne = tmp[i]
			}
		}
		if bigTime > okLen {
			volatilities = append(volatilities, StockVolatility{
				Name: sd.Symbol,
				Wave: lastOne,
			})
		}
	}

	// 根据波动性从高到低排序
	sort.Slice(volatilities, func(i, j int) bool {
		return volatilities[i].Wave > volatilities[j].Wave
	})

	return volatilities
}
