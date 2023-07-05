package protocol

type Request struct {
	Method string        `json:"Method"`
	Args   []interface{} `json:"Args"`
}
