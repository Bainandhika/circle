package tool

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v any) string {
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}
	return string(jsonByte)
}
