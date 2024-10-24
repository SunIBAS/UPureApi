# 纯 API

> 背景：多次反复复写 OKX 和 Bian 服务，过段时间就忘了如何使用，这里单纯按需封装 api

# 文档

1. [HttpUtils](./Core/HttpUtils/readme.md)
2. [BianHttpUtils](./Core/HttpUtils/BinaHttpUtils/readme.md)

- 文件目录

```text
D:.
│  go.mod
│  readme.md
│  
├─.idea
│      .gitignore
│      modules.xml
│      UPureApi.iml
│      workspace.xml
│
├─Core // 该文件夹下的内容如果是 go 文件则是简单的工具类
│  │  ts.go
│  │
│  └─HttpUtils  // 封装的请求类
│      │  Const.go
│      │  HttpUtils.go
│      │  readme.md
│      │
│      └─BinaHttpUtils  // bina 二次封装的 请求类
│          │  Api.go
│          │  httpUtils.go
│          │  Signalture.go
│          │
│          └─BinaApis   // bina 的 API 封装
│              │  Core.go   // 公用的枚举类
│              │
│              ├─KLine  // 对应 kline 目录下的接口
│              │      List.go
│              │
│              └─UsdsMarginedFuturesOrder // 对应 U 本位目录下的接口
│                      All.go
│                      const.go
│                      OpenOrder.go
│                      Order.go
│                      Query.go
│                      readme.md
│
├─OKex
│
├─Bina  以 bina api 组合出来的应用 
│
└─Test
    ├─Bina
    │      httpUtils_test.go
    │      sign_test.go
    │
    └─OKex

Bina bian api 接口实现
OKex ok api 接口实现
Core 工具类和核心类
```

### Bina 配置

- [API](https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/general-info)

- [Github Python Api](https://github.com/binance/binance-futures-connector-python/blob/main/binance/__init__.py)

1. 接口可能需要用户的 [API Key](https://www.binance.com/zh-CN/my/settings/api-management)，如何创建API-KEY请参考[这里](https://www.binance.com/zh-CN/support/faq/%E5%A6%82%E4%BD%95%E5%9C%A8%E5%B8%81%E5%AE%89%E5%88%9B%E5%BB%BAapi%E5%AF%86%E9%92%A5-360002502072?hl=zh-CN)
2. 本篇列出REST接口的baseurl https://fapi.binance.com
3. 所有接口的响应都是JSON格式 
4. 响应中如有数组，数组元素以时间升序排列，越早的数据越提前。 
5. 所有时间、时间戳均为UNIX时间，单位为毫秒 
6. 所有数据类型采用JAVA的数据类型定义"
