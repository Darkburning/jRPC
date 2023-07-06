package protocol

/*
	Request包含方法名和参数
	Response包含错误信息和返回值
*/

type Request struct {
	Method string        `json:"Method" yaml:"Method"`
	Args   []interface{} `json:"Args" yaml:"Args"`
}

type Response struct {
	Err     string        `json:"err" yaml:"err"`
	Replies []interface{} `json:"replies" yaml:"replies"`
}
