//Helper functions for converting sql.Row to
//a single object.

package db

import (
	"database/sql"
	"errors"
	"reflect"
)

func Row2Obj(row *sql.Row, objPtr interface{}) error {
	if row == nil {
		return errors.New("empty row")
	}

	objPtrVal := reflect.ValueOf(objPtr)
	objVal := reflect.Indirect(objPtrVal)

	if objPtrVal.Kind() == reflect.Invalid || objVal.Kind() == reflect.Invalid || objPtrVal.Kind() != reflect.Ptr ||
		objVal.Kind() == reflect.Ptr || objVal.Kind() == reflect.Chan || objVal.Kind() == reflect.Slice ||
		objVal.Kind() == reflect.Map || objVal.Kind() == reflect.Array || objVal.Kind() == reflect.Interface {
		return errors.New("not a pointer to supported types")
	}

	if objVal.Kind() == reflect.Struct {
		fieldCnt := objVal.NumField()
		dests := make([]interface{}, 0, fieldCnt)

		for i := 0; i < fieldCnt; i++ {
			fieldVal := objVal.Field(i)
			structField := objVal.Type().Field(i)

			if _, ok := structField.Tag.Lookup(ColTagName); !ok {
				continue
			}

			switch fieldVal.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var dest sql.NullInt64
				dests = append(dests, &dest)
			case reflect.Float32, reflect.Float64:
				var dest sql.NullFloat64
				dests = append(dests, &dest)
			case reflect.String:
				var dest sql.NullString
				dests = append(dests, &dest)
			case reflect.Bool:
				var dest sql.NullBool
				dests = append(dests, &dest)
			default:
				return errors.New("field type not supported")
			}
		}

		err := row.Scan(dests...)
		if err != nil {
			return err
		}

		for i, j := 0, 0; i < fieldCnt; i++ {
			fieldVal := objVal.Field(i)
			structField := objVal.Type().Field(i)

			if _, ok := structField.Tag.Lookup(ColTagName); !ok {
				continue
			}

			switch fieldVal.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val := dests[j].(*sql.NullInt64)
				if val.Valid {
					fieldVal.SetUint(uint64(val.Int64))
				}
				j++
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val := dests[j].(*sql.NullInt64)
				if val.Valid {
					fieldVal.SetInt(val.Int64)
				}
				j++
			case reflect.Float32, reflect.Float64:
				val := dests[j].(*sql.NullFloat64)
				if val.Valid {
					fieldVal.SetFloat(val.Float64)
				}
				j++
			case reflect.String:
				val := dests[j].(*sql.NullString)
				if val.Valid {
					fieldVal.SetString(val.String)
				}
				j++
			case reflect.Bool:
				val := dests[j].(*sql.NullBool)
				if val.Valid {
					fieldVal.SetBool(val.Bool)
				}
				j++
			}
		}
	} else {
		dests := make([]interface{}, 0, 1)

		switch objVal.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var dest sql.NullInt64
			dests = append(dests, &dest)
		case reflect.Float32, reflect.Float64:
			var dest sql.NullFloat64
			dests = append(dests, &dest)
		case reflect.String:
			var dest sql.NullString
			dests = append(dests, &dest)
		case reflect.Bool:
			var dest sql.NullBool
			dests = append(dests, &dest)
		default:
			return errors.New("object type not supported")
		}

		err := row.Scan(dests...)
		if err != nil {
			return err
		}

		switch objVal.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val := dests[0].(*sql.NullInt64)
			if val.Valid {
				objVal.SetUint(uint64(val.Int64))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := dests[0].(*sql.NullInt64)
			if val.Valid {
				objVal.SetInt(val.Int64)
			}
		case reflect.Float32, reflect.Float64:
			val := dests[0].(*sql.NullFloat64)
			if val.Valid {
				objVal.SetFloat(val.Float64)
			}
		case reflect.String:
			val := dests[0].(*sql.NullString)
			if val.Valid {
				objVal.SetString(val.String)
			}
		case reflect.Bool:
			val := dests[0].(*sql.NullBool)
			if val.Valid {
				objVal.SetBool(val.Bool)
			}
		}
	}

	objPtrVal.Elem().Set(objVal)

	return nil
}
