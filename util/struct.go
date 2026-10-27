package util

import "reflect"

func StructToMapReflect(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		field := reflect.TypeOf(data).Field(i).Name
		value := v.Field(i).Interface()
		result[field] = value
	}

	return result
}
