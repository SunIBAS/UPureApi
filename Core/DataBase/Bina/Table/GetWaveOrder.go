package Table

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GetWaveOrder struct {
	OrderId       string  `json:"symbol" gorm:"primaryKey"`
	Pair          string  `json:"pair"`
	Qty           float64 `json:"qty"`
	LongProfitId  float64 `json:"longProfitId"`
	ShortProfitId float64 `json:"shortProfitId"`
	LongStopId    float64 `json:"longStopId"`
	ShortStopId   float64 `json:"shortStopId"`
	LongProfit    float64 `json:"longProfit"`
	ShortProfit   float64 `json:"shortProfit"`
	LongStop      float64 `json:"longStop"`
	ShortStop     float64 `json:"shortStop"`
	StartTime     int64   `json:"startTime"`
}

// BeforeCreate https://gorm.io/zh_CN/docs/create.html
func (gwo *GetWaveOrder) BeforeCreate(tx *gorm.DB) (err error) {
	gwo.StartTime = time.Now().UnixMilli()
	var uuid uuid2.UUID
	uuid, err = uuid2.NewUUID()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("success")
	gwo.OrderId = uuid.String()
	return
}
