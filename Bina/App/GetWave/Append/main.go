package main

import (
	"UPureApi/Core/CoreUtils"
	"UPureApi/Core/DataBase/Bina"
	"UPureApi/Core/DataBase/Bina/Table"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis"
	"UPureApi/Core/HttpUtils/BinaHttpUtils/BinaApis/UsdtMarginedFutures/Market/UMFMarketKLine"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

type config struct {
	DbPath       string `json:"dbPath"`
	SaveLinePath string `json:"saveLinePath"`
	Proxy        struct {
		Porto string `json:"porto"`
		Host  string `json:"host"`
		Port  string `json:"port"`
	}
	Key struct {
		ApiKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
	}
	BaseUrl string `json:"base_url"`
	Log     bool   `json:"log"`
}

var conf config
var db *Bina.SQLite

func parseConfig() {
	bs, _ := ioutil.ReadFile(os.Args[1])
	json.Unmarshal(bs, &conf)
}

func main() {
	if len(os.Args) > 1 {
		parseConfig()
	} else {
		return
	}
	go runFetchLine()
	go runFetchDb()
	for {
	}
}

type job struct {
	Start int64                          `json:"start"`
	End   int64                          `json:"end"`
	Pair  string                         `json:"pair"`
	Order Table.GetWaveOrder             `json:"order"`
	Lines []UMFMarketKLine.KLineResponse `json:"lines"`
}

type jobQueue struct {
	job  []job
	lock sync.Mutex
}

func (jq *jobQueue) removeFirst() {
	jq.lock.Lock()
	defer jq.lock.Unlock()
	rj := queue.job[0]
	jq.job = jq.job[1:]
	fmt.Println(fmt.Sprintf("remove job [%s]", rj.Order.OrderId))
}
func (jq *jobQueue) add(j job) {
	jq.lock.Lock()
	defer jq.lock.Unlock()
	jq.job = append(jq.job, j)
	fmt.Println(fmt.Sprintf("add job [%s]", j.Order.OrderId))
}

var queue = jobQueue{
	job:  []job{},
	lock: sync.Mutex{},
}

func getJobFileByOrderId(orderId string) string {
	return path.Join(conf.SaveLinePath, orderId)
}

func runFetchLine() {
	server := BinaHttpUtils.NewBinaHttpUtilsFromConfig(BinaHttpUtils.BinaHttpUtilsConfig{
		Proxy:   conf.Proxy,
		Key:     conf.Key,
		BaseUrl: conf.BaseUrl,
		Log:     conf.Log,
	})
	for {
		if len(queue.job) == 0 {
			//time.Sleep(time.Millisecond * time.Duration(CoreUtils.Time1m))
		} else {
			j := queue.job[0]
			if fileExists(getJobFileByOrderId(j.Order.OrderId)) {
				// 如果是重复任务则丢弃
				queue.removeFirst()
				continue
			}
			err := _fetchKLineLoop(j, server)
			if err != nil {
				// 任何错误出现都导致改请求失败
			} else {
				if db != nil {
					db.GetDb().Delete(&Table.GetWaveOrder{OrderId: j.Order.OrderId})
				}
			}
		}
	}
}

// 如果 k 节点多于 1500 ，则需要多次请求
func _fetchKLineLoop(j job, server BinaHttpUtils.BinaHttpUtils) error {
	startTime := j.Start - j.Start%CoreUtils.Time1m
	endTime := j.End + CoreUtils.Time1m - j.End%CoreUtils.Time1m
	count := int((startTime - endTime) / CoreUtils.Time1m)
	times := [][]int64{}
	fetchTime := 0
	for {
		times = append(times, []int64{
			startTime, startTime + 1500*CoreUtils.Time1M,
		})
		count -= 1500
		fetchTime++
		if count < 0 {
			break
		}
		if fetchTime > 10 {
			// 不应该放这么久的
			break
		}
	}

	j.Lines = []UMFMarketKLine.KLineResponse{}

	for idx := 0; idx < fetchTime; idx++ {
		api := UMFMarketKLine.CreateKLineApi(
			UMFMarketKLine.KLineListApiParam{
				Symbol:    j.Pair,
				Interval:  BinaApis.Interval1m,
				StartTime: times[idx][0],
				EndTime:   times[idx][1],
				Limit:     1500,
			},
		)
		ret, err := server.Request(api)
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			fmt.Println(fmt.Sprintf("job success over [%s]", j.Order.OrderId))
			queue.removeFirst()
		}
		j.Lines = append(j.Lines, UMFMarketKLine.ParseKLineResponse(ret)...)
	}
	//fmt.Println(j)
	bs, err := json.Marshal(j)
	if err != nil {
		fmt.Println(fmt.Sprintf("json marshal err = [%s]", err.Error()))
		return err
	}
	err = ioutil.WriteFile(getJobFileByOrderId(j.Order.OrderId), bs, fs.ModeAppend)
	return err
}

func runFetchDb() {
	// ① 从数据库获取最后一个执行的订单数据
	db = &Bina.SQLite{
		Path: conf.DbPath,
	}
	db.Init()
	{
		var oldOrders = []Table.GetWaveOrder{}
		var orders []Table.GetWaveOrder
		// 降序获取所有数据
		db.GetDb().Order("start_time desc").Find(&orders)
		for _, order := range orders {
			orderFile := getJobFileByOrderId(order.OrderId)
			if fileExists(orderFile) {
				if db != nil {
					db.GetDb().Delete(&Table.GetWaveOrder{OrderId: order.OrderId})
				}
				continue
			}
			oldOrders = append(oldOrders, order)
		}
		fmt.Println(oldOrders)
		for i := 1; i < len(oldOrders); i++ {
			queue.add(job{
				Start: oldOrders[i].StartTime - CoreUtils.Time1H*12,
				End:   oldOrders[i-1].StartTime,
				Pair:  oldOrders[i].Pair,
				Order: oldOrders[i],
				Lines: nil,
			})
		}
		if len(oldOrders) > 0 {
			// 如果已经发起超过一天了
			if time.Now().UnixMilli()-oldOrders[0].StartTime > CoreUtils.Time1D {
				queue.add(job{
					Start: oldOrders[0].StartTime - CoreUtils.Time1H*12,
					End:   time.Now().UnixMilli(),
					Pair:  oldOrders[0].Pair,
					Order: oldOrders[0],
					Lines: nil,
				})
			}
		}
	}
	for {
		var orders []Table.GetWaveOrder
		db.GetDb().Order("start_time desc").Limit(2).Find(&orders)
		if len(orders) == 2 {
			var pOrder Table.GetWaveOrder
			var lOrder Table.GetWaveOrder
			if orders[0].StartTime < orders[1].StartTime {
				pOrder, lOrder = orders[0], orders[1]
			} else {
				pOrder, lOrder = orders[1], orders[0]
			}
			fromTs := pOrder.StartTime - CoreUtils.Time1H*12
			lastTs := lOrder.StartTime

			queue.add(job{
				Start: fromTs,
				End:   lastTs,
				Pair:  pOrder.Pair,
				Order: pOrder,
				Lines: nil,
			})
			time.Sleep(time.Duration(CoreUtils.Time1H) * time.Millisecond)
		} else {
			time.Sleep(time.Duration(CoreUtils.Time1m) * time.Millisecond)
		}
	}
}

// 判断文件是否存在
func fileExists(fp string) bool {
	fmt.Println(fmt.Sprintf("check file [%s]", fp))
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
