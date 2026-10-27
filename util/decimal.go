package util

import (
	"math"
	"reflect"
	"strconv"
)

func RoundToPrecision(value float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(value*factor) / factor
}

func GetFloat64(value interface{}) float64 {
	if value == nil {
		return 0.0
	}
	switch v := value.(type) {
	case float64:
		return v
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

func GetFieldValue(obj interface{}, fieldName string) float64 {
	v := reflect.ValueOf(obj)
	field := v.FieldByName(fieldName)
	if field.IsValid() && field.CanInterface() {
		return GetFloat64(field.Interface())
	}
	return 0.0
}
