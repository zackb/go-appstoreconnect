package encoding

import "encoding/json"

type JsonEncoder struct {
}

// NewJsonEncoder creates a new json encoder
func NewJsonEncoder() *JsonEncoder {
	return &JsonEncoder{}
}

func (e *JsonEncoder) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (e *JsonEncoder) Decode(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}
