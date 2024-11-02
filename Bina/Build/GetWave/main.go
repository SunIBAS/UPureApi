package main

import (
	"UPureApi/Bina/App/GetWave"
	"UPureApi/Bina/App/GetWave/GetWaveCore"
	"os"
)

func main() {
	argvLen := len(os.Args)
	if argvLen > 1 {
		opts := []GetWaveCore.Option{
			// 设置最大的杠杆倍数
			func(wave *GetWaveCore.GetWave) {
				wave.MaxMargin = 45
			},
		}
		GetWave.CreateApp(os.Args[1], opts...)
	}
}
