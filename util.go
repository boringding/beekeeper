package beekeeper

import (
	"fmt"
	"net/url"
	"reflect"
)

const ParamTagName = "param"

func ParseQueryStr(queryStr string, v interface{}) (err error) {
	vals, err := url.ParseQuery(queryStr)
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	reflectVal := reflect.ValueOf(v).Elem()
	reflectType := reflectVal.Type()

	for i := 0; i < reflectVal.NumField(); i++ {
		tagVal := reflectType.Field(i).Tag.Get(ParamTagName)
		val := vals.Get(tagVal)
		reflectVal.Field(i).SetString(val)
	}

	return
}
