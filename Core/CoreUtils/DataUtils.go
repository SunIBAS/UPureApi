package CoreUtils

import "strconv"

func StringToInt64(s string) int64 {
	i64, _ := strconv.ParseInt(s, 10, 64)
	return i64
}

func StringToFloat64(s string) float64 {
	f64, _ := strconv.ParseFloat(s, 10)
	return f64
}
