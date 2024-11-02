package main

import (
	"UPureApi/Core/CoreUtils/SingleTask"
	"fmt"
	"time"
)

type Coin struct {
	symbol string
	time   int
}

func main() {
	test2()
}

func test2() {
	sTask := SingleTask.NewProcessor(func(i interface{}) {
		coin := i.(Coin)
		fmt.Printf("[%s] %d start \r\n", coin.symbol, coin.time)
		// 2 秒执行结束
		time.Sleep(time.Second * 2)
		fmt.Printf("[%s] %d end \r\n", coin.symbol, coin.time)
	}, func(i interface{}) {
		coin := i.(Coin)
		fmt.Printf("[%s] %d drop \r\n", coin.symbol, coin.time)
	}, 100)
	coins := []Coin{
		{
			symbol: "BTC",
			time:   1,
		},
		{
			symbol: "ETH",
			time:   2,
		},
		{
			symbol: "SOL",
			time:   3,
		},
	}
	for _, coin := range coins {
		fmt.Printf("[%s] %d add \r\n", coin.symbol, coin.time)
		sTask.AddCoin(coin)
		time.Sleep(time.Millisecond * 500)
	}
	go sTask.Run()
	// [BTC] 1 add
	// [ETH] 2 add
	// [BTC] 1 drop
	// [SOL] 3 add
	// [ETH] 2 drop
	// [SOL] 3 start
	// [SOL] 3 end
	for {
	}
}

func test1() {
	sTask := SingleTask.NewProcessor(func(i interface{}) {
		coin := i.(Coin)
		fmt.Printf("[%s] %d start \r\n", coin.symbol, coin.time)
		// 2 秒执行结束
		time.Sleep(time.Second * 2)
		fmt.Printf("[%s] %d end \r\n", coin.symbol, coin.time)
	}, func(i interface{}) {
		coin := i.(Coin)
		fmt.Printf("[%s] %d drop \r\n", coin.symbol, coin.time)
	}, 100)
	go sTask.Run()
	coins := []Coin{
		{
			symbol: "BTC",
			time:   1,
		},
		{
			symbol: "ETH",
			time:   2,
		},
		{
			symbol: "SOL",
			time:   3,
		},
	}
	for _, coin := range coins {
		fmt.Printf("[%s] %d add \r\n", coin.symbol, coin.time)
		sTask.AddCoin(coin)
		time.Sleep(time.Millisecond * 500)
	}
	// [BTC] 1 add
	// [BTC] 1 start
	// [ETH] 2 add
	// [SOL] 3 add
	// [ETH] 2 drop
	// [BTC] 1 end
	// [SOL] 3 start
	// [SOL] 3 end
	for {
	}
}
