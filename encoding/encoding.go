package encoding

import "errors"

type Encoding int

const (
	Json Encoding = 1
	Tsv  Encoding = 2
	Csv  Encoding = 3
)

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
