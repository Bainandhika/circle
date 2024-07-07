package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ToJSONString(v any) string {
	json, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}

	return string(json)
}

func ToMap(data any /* should be a struct */) map[string]any {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]any)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if !field.IsZero() {
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag != "" {
				result[jsonTag] = field.Interface()
			}
		}
	}

	return result
}
