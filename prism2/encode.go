package prism2

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Marshaler interface {
	MarshalMessage() ([]byte, error)
}

func Marshal(v any) ([]byte, error) {
	if m, ok := v.(Marshaler); ok {
		return m.MarshalMessage()
	}

	return marshalMessage(v)
}

func marshalMessage(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return nil, fmt.Errorf("expected pointer, received invalid type %T", v)
	}

	fields, err := marshalFields(rv)
	if err != nil {
		return nil, err
	}

	return bytes.Join(fields, SeparatorField), nil
}

func marshalFields(val reflect.Value) ([][]byte, error) {
	var fields [][]byte

	switch val.Kind() {
	case reflect.Bool:
		return nil, errors.New("bool not supported")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValueInt64 := val.Int()
		fields = append(fields, stringToBytes(strconv.Itoa(int(fieldValueInt64))))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValueUint64 := val.Uint()
		fields = append(fields, stringToBytes(strconv.Itoa(int(fieldValueUint64))))
	case reflect.Float32, reflect.Float64:
		fieldValueFloat64 := val.Float()
		fields = append(fields, stringToBytes(strconv.FormatFloat(fieldValueFloat64, 'f', -1, 64)))
	case reflect.String:
		fieldValueString := val.String()
		fields = append(fields, stringToBytes(fieldValueString))
	case reflect.Slice, reflect.Array:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			fields = append(fields, val.Bytes())
		} else {
			for i := 0; i < val.Len(); i += 1 {
				newFields, err := marshalFields(val.Index(i))
				if err != nil {
					return nil, err
				}

				fields = append(fields, newFields...)
			}
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i += 1 {
			field := val.Field(i)
			nestedFields, err := marshalFields(field)
			if err != nil {
				return nil, err
			}

			fields = append(fields, nestedFields...)
		}
	case reflect.Ptr:
		unwrapped := val.Elem()

		if unwrapped.IsValid() {
			nestedFields, err := marshalFields(unwrapped)
			if err != nil {
				return nil, err
			}

			fields = append(fields, nestedFields...)
		}
	}

	return fields, nil
}
