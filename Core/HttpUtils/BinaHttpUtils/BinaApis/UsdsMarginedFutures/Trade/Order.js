// 挂多单止损
function longStopLoss(num, symbol, price, size) {
    var time = UnixNano() / 1000000;
    var param = null;
    // var param = "symbol=" + symbol + "&quantity=" + size.toString() + "&stopPrice=" + price.toString() + "&side=SELL" + "&type=STOP_MARKET" + "&timestamp=" + time.toString();                           // 单向持仓
    if (size != 0) {
        param = "symbol=" + symbol + "&quantity=" + size.toString() + "&stopPrice=" + price.toString() + "&side=SELL" + "&type=STOP_MARKET" + "&positionSide=LONG" + "&timestamp=" + time.toString();       // 双向持仓
    } else {
        param = "symbol=" + symbol + "&closePosition=true" + "&stopPrice=" + price.toString() + "&side=SELL" + "&type=STOP_MARKET" + "&positionSide=LONG" + "&timestamp=" + time.toString();                // 仓位止盈止损
    }
    var ret = exchanges[num].IO("api", "POST", "/fapi/v1/order", param);
    Log(exchanges[num].GetLabel(), ": 挂多单止损：", ret);
    return true;
}

// 挂多单止盈
function longTakeProfit(num, symbol, price, size) {
    var time = UnixNano() / 1000000;
    var param = null;
    if (size != 0) {
        param = "symbol=" + symbol + "&quantity=" + size.toString() + "&stopPrice=" + price.toString() + "&side=SELL" + "&type=TAKE_PROFIT_MARKET" + "&positionSide=LONG" + "&timestamp=" + time.toString();       // 双向持仓
    } else {
        param = "symbol=" + symbol + "&closePosition=true" + "&stopPrice=" + price.toString() + "&side=SELL" + "&type=TAKE_PROFIT_MARKET" + "&positionSide=LONG" + "&timestamp=" + time.toString();                // 仓位止盈止损
    }
    var ret = exchanges[num].IO("api", "POST", "/fapi/v1/order", param);
    Log(exchanges[num].GetLabel(), ": 挂多单止盈：", ret);
    return true;
}

// 挂空单止损
function shortStopLoss(num, symbol, price, size) {
    var time = UnixNano() / 1000000;
    var param = null;
    // var param = "symbol=" + symbol + "&quantity=" + size.toString() + "&stopPrice=" + price.toString() + "&side=BUY" + "&type=STOP_MARKET" + "&timestamp=" + time.toString();                          // 单向持仓
    if (size != 0) {
        param = "symbol=" + symbol + "&quantity=" + size.toString() + "&stopPrice=" + price.toString() + "&side=BUY" + "&type=STOP_MARKET" + "&positionSide=SHORT" + "&timestamp=" + time.toString();     // 双向持仓
    } else {
        param = "symbol=" + symbol + "&closePosition=true" + "&stopPrice=" + price.toString() + "&side=BUY" + "&type=STOP_MARKET" + "&positionSide=SHORT" + "&timestamp=" + time.toString();              // 仓位止盈止损
    }
    var ret = exchanges[num].IO("api", "POST", "/fapi/v1/order", param);
    Log(exchanges[num].GetLabel(), ": 挂空单止损：", ret);
    return true;
}

// 挂空单止盈
function shortTakeProfit(num, symbol, price, size) {
    var time = UnixNano() / 1000000;
    var param = null;
    if (size != 0) {
        param = "symbol=" + symbol + "&quantity=" + size.toString() + "&stopPrice=" + price.toString() + "&side=BUY" + "&type=TAKE_PROFIT_MARKET" + "&positionSide=SHORT" + "&timestamp=" + time.toString();     // 双向持仓
    } else {
        param = "symbol=" + symbol + "&closePosition=true" + "&stopPrice=" + price.toString() + "&side=BUY" + "&type=TAKE_PROFIT_MARKET" + "&positionSide=SHORT" + "&timestamp=" + time.toString();              // 仓位止盈止损
    }
    var ret = exchanges[num].IO("api", "POST", "/fapi/v1/order", param);
    Log(exchanges[num].GetLabel(), ": 挂空单止盈：", ret);
    return true;
}