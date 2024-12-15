package prism

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

type Unmarshaler interface {
	UnmarshalMessage([]byte) error
}

func Unmarshal(content []byte, v any) error {
	if u, ok := v.(Unmarshaler); ok {
		return u.UnmarshalMessage(content)
	}

	return unmarshalMessage(content, v)
}

func unmarshalMessage(content []byte, v any) error {
	val := reflect.Indirect(reflect.ValueOf(v))

	fieldsScanner := bufio.NewScanner(bytes.NewReader(content))
	fieldsScanner.Split(splitFieldsFunc)

	return unmarshalFields(val, fieldsScanner)
}

var errFieldCount = errors.New("field count mismatch")

func unmarshalFields(val reflect.Value, fields *bufio.Scanner) error {
	// If val implements Unmarshaler, call its UnmarshalMessage method
	if val.Type().Implements(reflect.TypeOf((*Unmarshaler)(nil)).Elem()) {
		fieldValue, err := fieldValueFromScanner(fields)
		if err != nil {
			return err
		}
		return val.Interface().(Unmarshaler).UnmarshalMessage([]byte(fieldValue))
	}

	switch val.Kind() {
	case reflect.Bool:
		fieldValue, err := fieldValueFromScanner(fields)
		if err != nil {
			return err
		}
		fieldValueInt64, err := strconv.ParseInt(fieldValue, 10, 64)
		if err != nil {
			return err
		}
		val.SetBool(fieldValueInt64 != 0)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue, err := fieldValueFromScanner(fields)
		if err != nil {
			return err
		}
		fieldValueInt64, err := strconv.ParseInt(fieldValue, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(fieldValueInt64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValue, err := fieldValueFromScanner(fields)
		if err != nil {
			return err
		}
		fieldValueUint64, err := strconv.ParseUint(fieldValue, 10, 64)
		if err != nil {
			return err
		}
		val.SetUint(fieldValueUint64)
	case reflect.Float32, reflect.Float64:
		fieldValue, err := fieldValueFromScanner(fields)
		if err != nil {
			return err
		}
		fieldValueFloat64, err := strconv.ParseFloat(fieldValue, 64)
		if err != nil {
			return err
		}
		val.SetFloat(fieldValueFloat64)
	case reflect.String:
		fieldValue, err := fieldValueFromScanner(fields)
		if err != nil {
			return err
		}
		val.SetString(fieldValue)
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			if !fields.Scan() {
				return nil
			}
			val.SetBytes(fields.Bytes())
			return nil
		}

		// TODO: evaluate usefulness of slice decoding
		// for {
		// 	// add empty element to the slice
		// 	newElem := reflect.New(val.Type().Elem()).Elem()
		// 	err := unmarshalFields(newElem, fields)
		// 	if err != nil {
		// 		if err == errFieldCount {
		// 			return nil
		// 		}
		// 		return err
		// 	}
		// 	val.Set(reflect.Append(val, newElem))
		// }
	case reflect.Array:
		for i := 0; i < val.Len(); i += 1 {
			err := unmarshalFields(val.Index(i), fields)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i += 1 {
			err := unmarshalFields(val.Field(i), fields)
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		unwrapped := val.Elem()
		if !unwrapped.IsValid() {
			newUnwrapped := reflect.New(val.Type().Elem())
			err := unmarshalFields(newUnwrapped, fields)
			if err != nil {
				return err
			}

			val.Set(newUnwrapped)
			return nil
		}

		err := unmarshalFields(unwrapped, fields)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitFieldsFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.IndexByte(data, SeparatorField1); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func fieldValueFromScanner(fields *bufio.Scanner) (string, error) {
	if !fields.Scan() {
		return "", errFieldCount
	}

	return fields.Text(), nil
}
