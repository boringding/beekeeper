//Utilities.

package beekeeper

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

const ParamTagName = "param"

//Parse a query string like "a=1&b=2&c=str" into an object.
//Parameter v should be a pointer to an object.
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

		switch reflectVal.Field(i).Type().Name() {
		case "uint", "uint32":
			uval, err := strconv.ParseUint(val, 10, 32)
			if err == nil {
				reflectVal.Field(i).SetUint(uval)
			}
		case "uint64":
			uval, err := strconv.ParseUint(val, 10, 64)
			if err == nil {
				reflectVal.Field(i).SetUint(uval)
			}
		case "int", "int32":
			ival, err := strconv.ParseInt(val, 10, 32)
			if err == nil {
				reflectVal.Field(i).SetInt(ival)
			}
		case "int64":
			ival, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				reflectVal.Field(i).SetInt(ival)
			}
		case "float32":
			fval, err := strconv.ParseFloat(val, 32)
			if err == nil {
				reflectVal.Field(i).SetFloat(fval)
			}
		case "float64":
			fval, err := strconv.ParseFloat(val, 64)
			if err == nil {
				reflectVal.Field(i).SetFloat(fval)
			}
		case "string":
			reflectVal.Field(i).SetString(val)
		case "bool":
			bval, err := strconv.ParseBool(val)
			if err == nil {
				reflectVal.Field(i).SetBool(bval)
			}
		}
	}

	return
}
