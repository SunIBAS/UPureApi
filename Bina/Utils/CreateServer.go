package Utils

import (
	"UPureApi/Core/HttpUtils/BinaHttpUtils"
	"encoding/json"
	"io/ioutil"
)

func CreateServe(configFile string) BinaHttpUtils.BinaHttpUtils {
	//configFile := "D:\\all_code\\UPureApi\\config\\Bina.json"
	var config BinaHttpUtils.BinaHttpUtilsConfig
	bs, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(bs, &config)
	return BinaHttpUtils.NewBinaHttpUtilsFromConfig(config)
}
