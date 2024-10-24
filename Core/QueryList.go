package Core

type QueryList struct {
	// 返回值 string 是返回信息
	// bool 是是否需要重试
	QueryFunc func(string) (data string, retry bool)
	// 重试次数
	RetryTimes int
}

// Query
// 这里是一个系列的请求，每次请求是异步的
func (queryList QueryList) Query(names []string) map[string]string {
	results := make(map[string]string, len(names))
	for _, name := range names {
		tryTime := queryList.RetryTimes
		for tryTime > 0 {
			// 执行每一个 query，等待其完成后进入下一个
			if ret, retry := queryList.QueryFunc(name); retry {
				tryTime--
			} else {
				results[name] = ret
				tryTime = 0
			}
		}
	}
	return results
}
