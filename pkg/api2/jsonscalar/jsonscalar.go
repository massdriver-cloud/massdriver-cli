package jsonscalar

import (
	"encoding/json"
	"errors"
)

func Marshal(v interface{}) ([]byte, error) {
	bytes, _ := json.Marshal(v)
	return json.Marshal(string(bytes))
}

func Unmarshal(data []byte, v interface{}) error {
	return errors.New("Unimplemented")
}
