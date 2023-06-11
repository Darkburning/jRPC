package protocol

/*
	Request包含方法名和参数
	Response包含错误信息和返回值
*/

type Request struct {
	Method string        `json:"method"`
	Args   []interface{} `json:"args"`
}

type Response struct {
	Err     string        `json:"err"`
	Replies []interface{} `json:"replies"`
}
