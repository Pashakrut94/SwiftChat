package handlers


import (
	"encoding/json"
)

func FormatResp(payload interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(payload, "", " ")
	}
	return json.Marshal(payload)
}
