package Bina

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	//$ echo -n "symbol=BTCUSDT&side=BUY&type=LIMIT&quantity=1&price=9000&timeInForce=GTC&recvWindow=5000&timestamp=1591702613943"
	//| openssl dgst -sha256 -hmac "2b5eb11e18796d12d88f13dc27dbbd02c2cc51ff7059765ed9821957d82bb4d9"
	//(stdin)= 3c661234138461fcc7a7d8746c6558c9842d4e10870d2ecbedf7777cad694af9

	//str := "symbol=BTCUSDT&side=BUY&type=LIMIT&quantity=1&price=9000&timeInForce=GTC&recvWindow=5000&timestamp=1591702613943"
	//ak := "2b5eb11e18796d12d88f13dc27dbbd02c2cc51ff7059765ed9821957d82bb4d9"
	//signalture := "3c661234138461fcc7a7d8746c6558c9842d4e10870d2ecbedf7777cad694af9"

	//str := "symbol=BANANUSDT&timestamp=1729352262029"
	//ak := "BSNdoFJX9D3qmnOw6iIfKawJuqozWQICQtbplecWo9er8x2xKpgkfQxS4NJBUg2d"
	//signalture := "e1e2e876f902aabfa3fc6a14ca311f371fe86f310375e42136d875abbe7f73e0"

	str := "marginType=CROSSED&symbol=ETHUSDT&timestamp=1729409908201"
	ak := "RsPHKzK9EN1V5w0xa5go1MQHZiBMIC12pAT3ObwkQ1uVa3ozqdtAPWiKCJa4SlU9"
	signalture := "667462bff23a11d7c929bad017521af4f2faa7409e331c5039f8a64b52e584ae"
	// symbol=BANANUSDT&timestamp=1729352262029&signature=e1e2e876f902aabfa3fc6a14ca311f371fe86f310375e42136d875abbe7f73e0

	signaltureNow := BinaHttpUtils.Sign(str, ak)
	if signalture != signaltureNow {
		fmt.Println(signaltureNow)
		panic("not same")
	}

	//  curl -H "X-MBX-APIKEY: vE3BDAL1gP1UaexugRLtteaAHg3UO8Nza20uexEuW1Kh3tVwQfFHdAiyjjY428o2" -X POST
	// '
	//https://fapi.binance.com/fapi/v1/order?
	//timestamp=1671090801999&recvWindow=9999999&symbol=BTCUSDT&side=SELL&type=MARKET&quantity=1.23&
	//signature=aap36wD5loVXizxvvPI3wz9Cjqwmb3KVbxoym0XeWG1jZq8umqrnSk8H8dkLQeySjgVY91Ufs%2BBGCW%2B4sZjQEpgAfjM76riNxjlD3coGGEsPsT2lG39R%2F1q72zpDs8pYcQ4A692NgHO1zXcgScTGgdkjp%2Brp2bcddKjyz5XBrBM%3D'
}
