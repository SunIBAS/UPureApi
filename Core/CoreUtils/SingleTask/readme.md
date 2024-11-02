# 设计思路

> 原本是给 getWave 中，持续执行 addCoin 方法时设置的任务逻辑  
> 逻辑：  
>   gw.bc.AddCoin("BTC", 1)  
>   gw.bc.AddCoin("ETH", 2)  
>   gw.bc.AddCoin("SOL", 3)  
>   当上面三个语句执行时，如果任务中没有没有 coin 需要处理，那么先对 BTC 进行处理（因为当前没任务，接到处理 BTC 任务，肯定是直接上）  
> 
>   如果这时候 BTC 没执行完，ETH 和 SOL 也被推了过来，
>   那么存在一种可能，ETH 的任务已经过时，只需要执行 SOL 的任务即可，
>   那么我们希望将 ETH 丢弃
>

```text
// [BTC] 1 add
// [BTC] 1 start
// [ETH] 2 add
// [SOL] 3 add
// [ETH] 2 drop
// [BTC] 1 end
// [SOL] 3 start
// [SOL] 3 end
```
