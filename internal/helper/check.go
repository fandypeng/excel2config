package helper

import (
	"reflect"
	"regexp"
)

func IsValidEmail(email string) bool {
	reg := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	if len(reg.Find([]byte(email))) > 0 {
		return true
	}
	return false
}

// 检查target是否在searchIn里存在，searchIn可以是数组或map
func Contains(searchIn interface{}, target interface{}) bool {
	arrayValue := reflect.ValueOf(searchIn)
	switch reflect.TypeOf(searchIn).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < arrayValue.Len(); i++ {
			if arrayValue.Index(i).Interface() == target {
				return true
			}
		}
	case reflect.Map:
		if arrayValue.MapIndex(reflect.ValueOf(target)).IsValid() {
			return true
		}
	}
	return false
}
