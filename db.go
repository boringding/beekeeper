//Database interface.

package beekeeper

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/boringding/beekeeper/db"
)

//Parameter container should be a pointer to slice or map.
func LoadQueryRows(rows *sql.Rows, container interface{}) error {
	ptrVal := reflect.ValueOf(container)
	val := reflect.Indirect(ptrVal)

	if ptrVal.Kind() == reflect.Invalid || val.Kind() == reflect.Invalid || ptrVal.Kind() != reflect.Ptr {
		return errors.New("container is not a pointer")
	}

	if val.Kind() == reflect.Slice {
		return db.Rows2Slice(rows, container)
	} else if val.Kind() == reflect.Map {
		return db.Rows2Map(rows, container)
	} else {
		return errors.New("container is not a pointer to slice or map")
	}
}
