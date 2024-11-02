package CoreUtils

import (
	"fmt"
	"math"
	"strings"
)

// Fix 函数将浮点数 n 四舍五入到指定的小数位 s
// fmt.Println(Fix(0.0123, 3)) // 输出: 0.012
// fmt.Println(Fix(0.0125, 3)) // 输出: 0.013
// fmt.Println(Fix(1.2345, 2)) // 输出: 1.23
// fmt.Println(Fix(1.2355, 2)) // 输出: 1.24
func Fix(n float64, s int) float64 {
	// 计算乘以 10 的 s 次方
	pow := math.Pow(10, float64(s))
	// 使用 math.Round 进行四舍五入
	return math.Round(n*pow) / pow
}

// GetPointLen 获取小数点位数
func GetPointLen(a float64) int {
	bStr := fmt.Sprintf("%.10f", a) // 保留足够的小数位以捕获精度
	bStrs := strings.Split(strings.Split(bStr, ".")[1], "")
	l := len(bStrs)
	for i := l - 1; i >= 0; i-- {
		if bStrs[i] != "0" {
			return i + 1
		}
	}
	return 0
}

// Align 将 a 对齐到 b 的最小位数
// a = 1.2121 , b = 0.001
// Align(a,b) = 1.212
// 将 a 对齐 b
func Align(a, b float64) float64 {
	// 计算 b 的小数位数
	bStr := fmt.Sprintf("%.10f", b) // 保留足够的小数位以捕获精度
	decimalIndex := len(bStr) - 1 - strings.Index(bStr, ".")
	if decimalIndex < 0 {
		decimalIndex = 0 // 如果没有小数点，默认小数位数为0
	}

	// 计算对齐的倍数
	pow := math.Pow(10, float64(decimalIndex))
	return math.Round(a*pow) / pow
}
