//Database interface.

package beekeeper

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/boringding/beekeeper/db"
)

//Parameter container should be a pointer to object.
func LoadQueryRow(row *sql.Row, container interface{}) error {
	return db.Row2Obj(row, container)
}

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

func QueryRow(handler *sql.DB, container interface{}, query string, args ...interface{}) error {
	if handler == nil {
		return errors.New("empty handler")
	}

	row := handler.QueryRow(query, args...)

	err := LoadQueryRow(row, container)
	if err != nil {
		return err
	}

	return nil
}

func QueryRows(handler *sql.DB, container interface{}, query string, args ...interface{}) error {
	if handler == nil {
		return errors.New("empty handler")
	}

	rows, err := handler.Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	err = LoadQueryRows(rows, container)
	if err != nil {
		return err
	}

	return nil
}
