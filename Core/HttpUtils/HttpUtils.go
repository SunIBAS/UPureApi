package HttpUtils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	transport *http.Transport
	client    *http.Client
}
type RequestOption func(r *Request)

// SetProxy 自定义代理配置
func SetProxy(porto, host, port string) RequestOption {
	return func(r *Request) {
		proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", porto, host, port))
		fmt.Println(proxyURL)
		if err != nil {
			fmt.Printf("解析代理服务器地址失败：%s\n", err)
			return
		}
		r.transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			//Proxy: http.ProxyFromEnvironment,
			// 跳过安全验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}

// SetProxyByEnv 使用系统环境的代理配置
func SetProxyByEnv() RequestOption {
	return func(r *Request) {
		r.transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
	}
}
func CreateRequest(options ...RequestOption) Request {
	r := Request{
		transport: &http.Transport{},
	}
	for _, i := range options {
		i(&r)
	}
	r.client = &http.Client{
		Transport: r.transport,
	}
	return r
}

var DefaultHeader = Map2Header(map[string]string{
	"Accept":       "application/json",
	"Content-Type": "application/json",
})

func Map2Header(header map[string]string) http.Header {
	hd := make(http.Header)
	for k, v := range header {
		hd.Add(k, v)
	}
	return hd
}

func (r *Request) ToRequest(method HttpMethod, urlString string, headers http.Header, reqBody io.Reader) (string, error) {
	// 创建HTTP请求
	req, err := http.NewRequest(method.MethodName(), urlString, reqBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// 将自定义头设置到请求中
	req.Header = headers
	// 发送HTTP请求
	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
func (r *Request) Get(urlString string, headers http.Header) error {
	//client := &http.Client{
	//	Transport: r.transport,
	//	Timeout:   time.Second * 5,
	//}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// 将自定义头设置到请求中
	req.Header = headers
	//req.URL.Query()
	// 发送HTTP请求
	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

type P2SJSON int

const (
	ToJson   P2SJSON = 0
	ToString P2SJSON = 1
)

// ToJson 转 json 格式 {"a": "b", "c": "d", ...}
// ToString 转 urlString 格式 "a=b&c=d&..."
func params2string(params map[string]string, method P2SJSON) string {
	var paramsStr string
	if method == ToString {
		value := url.Values{}
		for k, v := range params {
			value.Add(k, v)
		}
		paramsStr = value.Encode()

	} else if method == ToJson {
		data, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Error marshaling json:", err)
			panic(err)
		}
		paramsStr = string(data)
	}
	return paramsStr
}
func Params2string(params map[string]string, method P2SJSON) string {
	return params2string(params, method)
}
