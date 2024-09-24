package utils

import "encoding/json"

func JsonMarshal(v any) string {
	bs, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bs)
}
