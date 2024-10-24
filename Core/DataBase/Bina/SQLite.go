package Bina

import (
	"UPureApi/Core/DataBase/Bina/Table"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type SQLiteConfig struct {
	Path string
}
type SQLite struct {
	Path string `json:"path"`
	db   *gorm.DB
}

func (sq *SQLite) Init() error {
	var err error = nil
	// 连接 SQLite 数据库
	sq.db, err = gorm.Open(sqlite.Open(sq.Path), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
		return err
	}

	// 自动迁移（创建表）
	sq.db.AutoMigrate(&Table.ExchangeInfoSymbol{})
	sq.db.AutoMigrate(&Table.Config{})
	sq.db.AutoMigrate(&Table.KLine{})
	sq.db.AutoMigrate(&Table.Brackets{})

	return err
}

func (sq *SQLite) GetDb() *gorm.DB {
	return sq.db
}
