package Bina

import (
	"UPureApi/Core/DataBase/Bina"
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

func createServe() BinaHttpUtils.BinaHttpUtils {
	configFile := "D:\\all_code\\UPureApi\\config\\Bina.json"
	var config BinaHttpUtils.BinaHttpUtilsConfig
	bs, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(bs, &config)
	return BinaHttpUtils.NewBinaHttpUtilsFromConfig(config)
}

// dbName 为了方便在不同测试中共享数据库，这里使用 dbName 便于快捷使用相同数据库
func createDb(dbName string) Bina.SQLite {
	return Bina.SQLite{
		Path: path.Join("D:\\all_code\\UPureApi\\Test\\Bina\\db", fmt.Sprintf("%s.db", dbName)),
	}
}
