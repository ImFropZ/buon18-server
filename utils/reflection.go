package utils

import "reflect"

func IsAllFieldsNil(v interface{}) bool {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsNil() {
			return false
		}
	}
	return true
}
