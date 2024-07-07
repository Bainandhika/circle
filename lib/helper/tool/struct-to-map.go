package tool

import (
	"encoding/json"
	"reflect"
)

func StructToMap(s any, resultInJson bool) map[string]any {
	if resultInJson {
		dataByte, _ := json.Marshal(s)

		var result map[string]any
		_ = json.Unmarshal(dataByte, &result)
		return result
	}

	result := make(map[string]any)
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return result
}
