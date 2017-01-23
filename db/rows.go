//Helper functions for converting sql.Rows to
//object slice or map.

package db

import (
	"database/sql"
	"errors"
	"reflect"
)

const (
	ColTagName = "col"
	IdTagName  = "id"
)

func Rows2Slice(rows *sql.Rows, slicePtr interface{}) error {
	if rows == nil {
		return errors.New("empty rows")
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	slicePtrVal := reflect.ValueOf(slicePtr)
	sliceVal := reflect.Indirect(slicePtrVal)

	if slicePtrVal.Kind() == reflect.Invalid || sliceVal.Kind() == reflect.Invalid ||
		slicePtrVal.Kind() != reflect.Ptr || sliceVal.Kind() != reflect.Slice {
		return errors.New("not a pointer to slice")
	}

	elemType := sliceVal.Type().Elem()

	if elemType.Kind() == reflect.Ptr || elemType.Kind() == reflect.Chan || elemType.Kind() == reflect.Slice ||
		elemType.Kind() == reflect.Map || elemType.Kind() == reflect.Array || elemType.Kind() == reflect.Interface {
		return errors.New("element type not supported")
	}

	newElem := reflect.New(elemType).Elem()
	dests := make([]interface{}, 0, len(cols))

	if elemType.Kind() == reflect.Struct {
		for i := 0; i < newElem.NumField(); i++ {
			fieldVal := newElem.Field(i)
			structField := newElem.Type().Field(i)

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

		if len(dests) != len(cols) {
			return errors.New("can not match object to columns")
		}
	} else {
		switch newElem.Kind() {
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
			return errors.New("element type not supported")
		}
	}

	for rows.Next() {
		err = rows.Scan(dests...)
		if err != nil {
			return err
		}

		if elemType.Kind() == reflect.Struct {
			for i, j := 0, 0; i < newElem.NumField(); i++ {
				fieldVal := newElem.Field(i)
				structField := newElem.Type().Field(i)

				if _, ok := structField.Tag.Lookup(ColTagName); !ok {
					continue
				}

				switch fieldVal.Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val := dests[j].(*sql.NullInt64)
					if val.Valid {
						fieldVal.SetUint(uint64(val.Int64))
					} else {
						fieldVal.SetUint(uint64(0))
					}
					j++
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val := dests[j].(*sql.NullInt64)
					if val.Valid {
						fieldVal.SetInt(val.Int64)
					} else {
						fieldVal.SetInt(0)
					}
					j++
				case reflect.Float32, reflect.Float64:
					val := dests[j].(*sql.NullFloat64)
					if val.Valid {
						fieldVal.SetFloat(val.Float64)
					} else {
						fieldVal.SetFloat(0.0)
					}
					j++
				case reflect.String:
					val := dests[j].(*sql.NullString)
					if val.Valid {
						fieldVal.SetString(val.String)
					} else {
						fieldVal.SetString("")
					}
					j++
				case reflect.Bool:
					val := dests[j].(*sql.NullBool)
					if val.Valid {
						fieldVal.SetBool(val.Bool)
					} else {
						fieldVal.SetBool(false)
					}
					j++
				}
			}
		} else {
			switch newElem.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val := dests[0].(*sql.NullInt64)
				if val.Valid {
					newElem.SetUint(uint64(val.Int64))
				} else {
					newElem.SetUint(uint64(0))
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val := dests[0].(*sql.NullInt64)
				if val.Valid {
					newElem.SetInt(val.Int64)
				} else {
					newElem.SetInt(0)
				}
			case reflect.Float32, reflect.Float64:
				val := dests[0].(*sql.NullFloat64)
				if val.Valid {
					newElem.SetFloat(val.Float64)
				} else {
					newElem.SetFloat(0.0)
				}
			case reflect.String:
				val := dests[0].(*sql.NullString)
				if val.Valid {
					newElem.SetString(val.String)
				} else {
					newElem.SetString("")
				}
			case reflect.Bool:
				val := dests[0].(*sql.NullBool)
				if val.Valid {
					newElem.SetBool(val.Bool)
				} else {
					newElem.SetBool(false)
				}
			}
		}

		sliceVal = reflect.Append(sliceVal, newElem)
	}

	slicePtrVal.Elem().Set(sliceVal)

	return nil
}

func Rows2Map(rows *sql.Rows, mapPtr interface{}) error {
	if rows == nil {
		return errors.New("empty rows")
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	mapPtrVal := reflect.ValueOf(mapPtr)
	mapVal := reflect.Indirect(mapPtrVal)

	if mapPtrVal.Kind() == reflect.Invalid || mapVal.Kind() == reflect.Invalid ||
		mapPtrVal.Kind() != reflect.Ptr || mapVal.Kind() != reflect.Map {
		return errors.New("not a pointer to map")
	}

	keyType := mapVal.Type().Key()
	elemType := mapVal.Type().Elem()

	if keyType.Kind() == reflect.Ptr || keyType.Kind() == reflect.Chan || keyType.Kind() == reflect.Slice ||
		keyType.Kind() == reflect.Map || keyType.Kind() == reflect.Array || keyType.Kind() == reflect.Interface ||
		elemType.Kind() == reflect.Ptr || elemType.Kind() == reflect.Chan || elemType.Kind() == reflect.Slice ||
		elemType.Kind() == reflect.Map || elemType.Kind() == reflect.Array || elemType.Kind() == reflect.Interface {
		return errors.New("key or element type not supported")
	}

	newKey := reflect.New(keyType).Elem()
	newElem := reflect.New(elemType).Elem()
	dests := make([]interface{}, 0, len(cols))

	if elemType.Kind() == reflect.Struct {
		for i := 0; i < newElem.NumField(); i++ {
			fieldVal := newElem.Field(i)
			structField := newElem.Type().Field(i)

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

		if len(dests) != len(cols) {
			return errors.New("can not match object to columns")
		}
	} else {
		switch newElem.Kind() {
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
			return errors.New("element type not supported")
		}
	}

	for rows.Next() {
		err = rows.Scan(dests...)
		if err != nil {
			return err
		}

		if elemType.Kind() == reflect.Struct {
			for i, j := 0, 0; i < newElem.NumField(); i++ {
				fieldVal := newElem.Field(i)
				structField := newElem.Type().Field(i)
				isId := false

				if _, ok := structField.Tag.Lookup(ColTagName); !ok {
					continue
				}

				if _, ok := structField.Tag.Lookup(IdTagName); ok {
					isId = true
				}

				switch fieldVal.Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val := dests[j].(*sql.NullInt64)
					if val.Valid {
						fieldVal.SetUint(uint64(val.Int64))

						if isId == true {
							newKey.SetUint(uint64(val.Int64))
						}
					} else {
						fieldVal.SetUint(uint64(0))

						if isId == true {
							newKey.SetUint(uint64(0))
						}
					}
					j++
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val := dests[j].(*sql.NullInt64)
					if val.Valid {
						fieldVal.SetInt(val.Int64)

						if isId == true {
							newKey.SetInt(val.Int64)
						}
					} else {
						fieldVal.SetInt(0)

						if isId == true {
							newKey.SetInt(0)
						}
					}
					j++
				case reflect.Float32, reflect.Float64:
					val := dests[j].(*sql.NullFloat64)
					if val.Valid {
						fieldVal.SetFloat(val.Float64)

						if isId == true {
							newKey.SetFloat(val.Float64)
						}
					} else {
						fieldVal.SetFloat(0.0)

						if isId == true {
							newKey.SetFloat(0.0)
						}
					}
					j++
				case reflect.String:
					val := dests[j].(*sql.NullString)
					if val.Valid {
						fieldVal.SetString(val.String)

						if isId == true {
							newKey.SetString(val.String)
						}
					} else {
						fieldVal.SetString("")

						if isId == true {
							newKey.SetString("")
						}
					}
					j++
				case reflect.Bool:
					val := dests[j].(*sql.NullBool)
					if val.Valid {
						fieldVal.SetBool(val.Bool)

						if isId == true {
							newKey.SetBool(val.Bool)
						}
					} else {
						fieldVal.SetBool(false)

						if isId == true {
							newKey.SetBool(false)
						}
					}
					j++
				}
			}
		} else {
			switch newElem.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val := dests[0].(*sql.NullInt64)
				if val.Valid {
					newElem.SetUint(uint64(val.Int64))
					newKey.SetUint(uint64(val.Int64))
				} else {
					newElem.SetUint(uint64(0))
					newKey.SetUint(uint64(0))
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val := dests[0].(*sql.NullInt64)
				if val.Valid {
					newElem.SetInt(val.Int64)
					newKey.SetInt(val.Int64)
				} else {
					newElem.SetInt(0)
					newKey.SetInt(0)
				}
			case reflect.Float32, reflect.Float64:
				val := dests[0].(*sql.NullFloat64)
				if val.Valid {
					newElem.SetFloat(val.Float64)
					newKey.SetFloat(val.Float64)
				} else {
					newElem.SetFloat(0.0)
					newKey.SetFloat(0.0)
				}
			case reflect.String:
				val := dests[0].(*sql.NullString)
				if val.Valid {
					newElem.SetString(val.String)
					newKey.SetString(val.String)
				} else {
					newElem.SetString("")
					newKey.SetString("")
				}
			case reflect.Bool:
				val := dests[0].(*sql.NullBool)
				if val.Valid {
					newElem.SetBool(val.Bool)
					newKey.SetBool(val.Bool)
				} else {
					newElem.SetBool(false)
					newKey.SetBool(false)
				}
			}
		}

		mapVal.SetMapIndex(newKey, newElem)
	}

	mapPtrVal.Elem().Set(mapVal)

	return nil
}
