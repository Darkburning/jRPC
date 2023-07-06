package serializer

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type YamlSerializer struct {
}

var _ Serializer = (*YamlSerializer)(nil)

func (y *YamlSerializer) Marshal(message interface{}) ([]byte, error) {
	if message == nil {
		return []byte{}, nil
	}
	return yaml.Marshal(message)
}

func (y *YamlSerializer) Unmarshal(data []byte, message interface{}) error {
	if len(data) == 0 {
		return errors.New("message empty")
	}
	return yaml.Unmarshal(data, message)
}
