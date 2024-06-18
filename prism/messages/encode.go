package messages

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/emilekm/go-prbf2/prism"
)

func (m baseMessage) Content() []byte {
	content, err := EncodeContent(m)
	if err != nil {
		panic(err)
	}

	return content
}

func EncodeContent(msg any) ([]byte, error) {
	val := reflect.ValueOf(msg)

	fields, err := encode(val)
	if err != nil {
		return nil, err
	}

	return bytes.Join(fields, prism.SeparatorField), nil
}

func encode(val reflect.Value) ([][]byte, error) {
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

func stringToBytes(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}
