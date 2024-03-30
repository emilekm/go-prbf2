package prism

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

type MessageEncoder interface {
	Encode() ([]byte, error)
}

type FieldsEncoder interface {
	EncodeFields() ([][]byte, error)
}

func Encode(msg Message) ([]byte, error) {
	var fields [][]byte
	var err error

	if encoder, ok := msg.(MessageEncoder); ok {
		return encoder.Encode()
	}

	val := reflect.ValueOf(msg)

	fields, err = encode(val)
	if err != nil {
		return nil, err
	}

	return bytes.Join(
		[][]byte{
			SeparatorStart,
			stringToBytes(string(msg.Subject())),
			SeparatorSubject,
			bytes.Join(fields, SeparatorField),
			SeparatorEnd,
			SeparatorNull,
		},
		[]byte{},
	), nil
}

func encode(val reflect.Value) ([][]byte, error) {
	var fields [][]byte

	if val.Type().Implements(reflect.TypeOf((*FieldsEncoder)(nil)).Elem()) && val.CanInterface() {
		dec := val.Interface().(FieldsEncoder)
		fields, err := dec.EncodeFields()
		if err != nil {
			return nil, err
		}
		return fields, nil
	}

	switch val.Kind() {
	case reflect.Bool:
		return nil, errors.New("bool not supported")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValueInt64 := val.Int()
		fields = append(fields, stringToBytes(strconv.Itoa(int(fieldValueInt64))))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValueUint64 := val.Uint()
		fields = append(fields, stringToBytes(strconv.Itoa(int(fieldValueUint64))))
	case reflect.String:
		fieldValueString := val.String()
		fields = append(fields, stringToBytes(fieldValueString))
	case reflect.Slice, reflect.Array:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			fields = append(fields, val.Bytes())
		} else {
			for i := 0; i < val.Len(); i += 1 {
				newFields, err := encode(val.Index(i))
				if err != nil {
					return nil, err
				}

				fields = append(fields, newFields...)
			}
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i += 1 {
			field := val.Field(i)
			nestedFields, err := encode(field)
			if err != nil {
				return nil, err
			}

			fields = append(fields, nestedFields...)
		}
	case reflect.Ptr:
		unwrapped := val.Elem()

		if unwrapped.IsValid() {
			nestedFields, err := encode(unwrapped)
			if err != nil {
				return nil, err
			}

			fields = append(fields, nestedFields...)
		}
	}

	return fields, nil
}
