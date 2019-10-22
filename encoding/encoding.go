package encoding

import "errors"

// Encoder defines a type for converting to and from different encodings
type Encoder interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type Encodable interface {
	GetHeader() []string
	Values() []string
}

// Encoding defines an available encoding format
type Encoding int

const (
	Json Encoding = 1
	Tsv  Encoding = 2
	Csv  Encoding = 3
)

func NewEncoder(e Encoding) Encoder {
	switch e {
	case Json:
		return NewJsonEncoder()
	default:
		panic("I dont know how to create an encoder: " + e.String())
	}
}

func (e *Encoding) Set(s string) error {
	switch s {
	case "csv":
		*e = Csv
	case "tsv":
		*e = Tsv
	case "json":
		*e = Json
	default:
		return errors.New("invalid output format: " + s)
	}
	return nil
}

func (e *Encoding) String() string {
	switch *e {
	case Csv:
		return "csv"
	case Tsv:
		return "tsv"
	case Json:
		return "json"
	}

	return "?"
}
