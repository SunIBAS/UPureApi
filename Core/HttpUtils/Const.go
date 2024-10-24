package HttpUtils

type HttpMethod int

const (
	HttpGet  = HttpMethod(1)
	HttpPost = HttpMethod(2)
)

type Api struct {
	Method HttpMethod `json:"method,omitempty"`
	Path   string     `json:"path,omitempty"`
	Auth   bool       `json:"auth,omitempty"`
}

func (hm HttpMethod) MethodName() string {
	if hm == HttpGet {
		return "GET"
	} else if hm == HttpPost {
		return "POST"
	} else {
		return ""
	}
}

type ApiCode string

const (
	ApiSuccess = ApiCode("0")
)

type ResendBody struct {
	Api    Api
	Params string
}
