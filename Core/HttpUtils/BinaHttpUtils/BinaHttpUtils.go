package BinaHttpUtils

import (
	"UPureApi/Core/HttpUtils/HttpUtilsCore"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BinaHttpUtilsConfig struct {
	Proxy struct {
		Porto string `json:"porto"`
		Host  string `json:"host"`
		Port  string `json:"port"`
	} `json:"proxy"`
	Key struct {
		ApiKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
	} `json:"key"`
	BaseUrl string `json:"base_url"`
	Log     bool   `json:"log"`
}

type BinaHttpUtils struct {
	ApiKey    string
	SecretKey string
	BaseUrl   string
	request   HttpUtilsCore.Request
	Log       bool
}

func NewBinaHttpUtilsFromConfig(config BinaHttpUtilsConfig) BinaHttpUtils {
	return NewBinaHttpUtils(
		HttpUtilsCore.SetProxy(config.Proxy.Porto, config.Proxy.Host, config.Proxy.Port),
		config.Key.ApiKey,
		config.Key.SecretKey,
		config.BaseUrl,
		config.Log,
	)
}
func NewBinaHttpUtils(option HttpUtilsCore.RequestOption, ApiKey, SecretKey, BaseUrl string, log bool) BinaHttpUtils {
	return BinaHttpUtils{
		request:   HttpUtilsCore.CreateRequest(option),
		SecretKey: SecretKey,
		ApiKey:    ApiKey,
		BaseUrl:   BaseUrl,
		Log:       log,
	}
}

func (binaHttpUtils *BinaHttpUtils) sign(str string) string {
	signResult := Sign(str, binaHttpUtils.SecretKey)
	if binaHttpUtils.Log {
		fmt.Printf(`str = [%s]
key = [%s]
result = [%s]
`, str, binaHttpUtils.SecretKey, signResult)
	}
	return signResult
}

func (binaHttpUtils *BinaHttpUtils) Request(api Api) (string, error) {
	return binaHttpUtils.RequestL(api, false)
}
func (binaHttpUtils *BinaHttpUtils) RequestL(api Api, log bool) (string, error) {
	queryStringMap := api.QueryParams.ToMap()
	//bodyMap := api.QueryParams.ToMap()
	//queryStringMap["recvWindow"] = "6000000"
	//queryStringMap["timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	timestamp := fmt.Sprintf("timestamp=%s", strconv.FormatInt(time.Now().UnixMilli(), 10))
	if api.NoTimeStamp {
		timestamp = ""
	}
	queryString := queryStringMap.ToString()
	//url := fmt.Sprintf("%s?%s", api.Path, queryString)
	//if api.Header == nil {
	//	api.Header = HttpUtils.DefaultHeader
	//}
	api.Header = http.Header{}
	if queryString != "" {
		queryString = fmt.Sprintf("%s&%s", queryString, timestamp)
	} else {
		queryString = timestamp
	}
	//url := fmt.Sprintf("%s?%s", api.Path, queryString)
	urlPart := []string{
		api.Path,
	}
	if queryString != "" {
		urlPart = append(urlPart, queryString)
	}
	if api.Sign {
		//reqBody["signature"] = binaHttpUtils.sign(reqBodyStr)
		urlPart = append(urlPart, fmt.Sprintf("signature=%s", binaHttpUtils.sign(queryString)))
		//url = fmt.Sprintf("%s&signature=%s", url, binaHttpUtils.sign(queryString))
		api.Header.Add("X-MBX-APIKEY", binaHttpUtils.ApiKey)
	}
	url := ""
	if len(urlPart) == 1 {
		url = urlPart[0]
	} else {
		url = fmt.Sprintf("%s?%s", urlPart[0], strings.Join(urlPart[1:], "&"))
	}

	if log {
		fmt.Println(fmt.Sprintf("[%s] %s", api.HttpMethod.MethodName(), url))
	}

	if api.HttpMethod == HttpUtilsCore.HttpPost {
		return binaHttpUtils.Post(
			url,
			api.Header,
			"",
			//bodyMap.ToString(),
		)
	} else if api.HttpMethod == HttpUtilsCore.HttpGet {
		return binaHttpUtils.Get(
			url,
			api.Header,
		)
	} else if api.HttpMethod == HttpUtilsCore.HttpDelete {
		return binaHttpUtils.Delete(
			url,
			api.Header,
		)
	}
	return "", errors.New(fmt.Sprintf("[BinaHttpUtils.Request] method = %d not found.", api.HttpMethod))
}
func (binaHttpUtils *BinaHttpUtils) Post(Path string, header http.Header, params string) (string, error) {
	if binaHttpUtils.Log {
		fmt.Printf("method = POST\r\nurl = %s\r\nheader = %s\r\n", fmt.Sprintf("%s%s", binaHttpUtils.BaseUrl, Path), header)
	}
	host := binaHttpUtils.BaseUrl
	if Path[0:4] == "http" {
		host = ""
	}
	return binaHttpUtils.request.ToRequest(
		HttpUtilsCore.HttpPost,
		fmt.Sprintf("%s%s", host, Path),
		header,
		nil,
	)
}
func (binaHttpUtils *BinaHttpUtils) Get(Path string, header http.Header) (string, error) {
	if binaHttpUtils.Log {
		fmt.Printf("method = POST\r\nurl = %s\r\nheader = %s\r\n", fmt.Sprintf("%s%s", binaHttpUtils.BaseUrl, Path), header)
	}
	host := binaHttpUtils.BaseUrl
	if Path[0:4] == "http" {
		host = ""
	}
	return binaHttpUtils.request.ToRequest(
		HttpUtilsCore.HttpGet,
		fmt.Sprintf("%s%s", host, Path),
		header,
		nil,
	)
}
func (binaHttpUtils *BinaHttpUtils) Delete(Path string, header http.Header) (string, error) {
	if binaHttpUtils.Log {
		fmt.Printf("method = POST\r\nurl = %s\r\nheader = %s\r\n", fmt.Sprintf("%s%s", binaHttpUtils.BaseUrl, Path), header)
	}
	host := binaHttpUtils.BaseUrl
	if Path[0:4] == "http" {
		host = ""
	}
	return binaHttpUtils.request.ToRequest(
		HttpUtilsCore.HttpDelete,
		fmt.Sprintf("%s%s", host, Path),
		header,
		nil,
	)
}
