package db

import (
	"UPureApi/Core/DataBase/Bina"
	"UPureApi/Core/DataBase/Bina/Table"
	"testing"
)

func TestDb(t *testing.T) {
	db := Bina.SQLite{
		Path: "D:\\all_code\\UPureApi\\config\\gw.db",
	}
	db.Init()
	db.GetDb().Create(&Table.GetWaveOrder{
		OrderId:       "",
		Pair:          "",
		Qty:           0,
		LongProfitId:  0,
		ShortProfitId: 0,
		LongStopId:    0,
		ShortStopId:   0,
		LongProfit:    0,
		ShortProfit:   0,
		LongStop:      0,
		ShortStop:     0,
		StartTime:     0,
	})
}
