package serializer

import (
	"encoding/json"
	"errors"
)

type JsonSerializer struct {
}

var _ Serializer = (*JsonSerializer)(nil)

func (j *JsonSerializer) Marshal(message interface{}) ([]byte, error) {
	if message == nil {
		return []byte{}, nil
	}
	return json.Marshal(message)
}

func (j *JsonSerializer) Unmarshal(data []byte, message interface{}) error {
	if len(data) == 0 {
		return errors.New("message empty")
	}
	return json.Unmarshal(data, message)
}