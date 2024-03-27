package prism

import (
	"errors"
	"reflect"
	"strconv"
)

func UnmarshalInto[T any](msg Message, into *T) error {
	val := reflect.ValueOf(into)
	val = val.Elem()

	fields := msg.Fields()

	numFields := val.NumField()
	if numFields > len(fields) {
		return errors.New("not enough fields")
	}

	for i := 0; i < val.NumField(); i += 1 {
		field := val.Field(i)
		fieldType := field.Type()
		fieldValue := fields[i]

		switch fieldType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// take fieldValue as int64
			fieldValueInt64, err := strconv.ParseInt(BytesToString(fieldValue), 10, 64)
			if err != nil {
				return err
			}
			// set field value
			field.SetInt(fieldValueInt64)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// take fieldValue as uint64
			fieldValueUint64, err := strconv.ParseUint(BytesToString(fieldValue), 10, 64)
			if err != nil {
				return err
			}

			field.SetUint(fieldValueUint64)
		case reflect.String:
			fieldValueString := BytesToString(fieldValue)
			field.SetString(fieldValueString)
		case reflect.Slice:
			if fieldType.Elem().Kind() == reflect.Uint8 {
				field.SetBytes(fieldValue)
			}
		}
	}

	return nil
}
