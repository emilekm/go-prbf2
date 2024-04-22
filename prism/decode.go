package prism

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

type ContentDecoder interface {
	Decode(content []byte) error
}

func decodeData(data []byte) (Subject, []byte, error) {
	// Make sure we arent reading data after the end separator
	// or with it
	data, _, _ = bytes.Cut(data, SeparatorEnd)
	// Skipping data before the start separator
	// It might be empty spaces
	_, data, found := bytes.Cut(data, SeparatorStart)
	if !found {
		return "", nil, errors.New("missing start separator")
	}

	subject, content, found := bytes.Cut(data, SeparatorSubject)
	if !found {
		return "", nil, errors.New("missing subject separator")
	}

	return Subject(subject), content, nil
}

func DecodeRaw(data []byte) (*RawMessage, error) {
	subject, content, err := decodeData(data)
	if err != nil {
		return nil, err
	}

	return &RawMessage{
		subject: subject,
		content: content,
	}, nil
}

func Decode[T Message](data []byte) (*T, error) {
	subject, content, err := decodeData(data)
	if err != nil {
		return nil, err
	}

	var msg T
	if msg.Subject() != subject {
		return nil, errors.New("subject mismatch")
	}

	err = DecodeContent(content, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, err
}

func DecodeContent(content []byte, into any) error {
	if decoder, ok := into.(ContentDecoder); ok {
		return decoder.Decode(content)
	}

	val := reflect.Indirect(reflect.ValueOf(into))

	fieldsScanner := bufio.NewScanner(bytes.NewReader(content))
	fieldsScanner.Split(splitFieldsFunc)

	return decode(val, fieldsScanner)
}

var errFieldCount = errors.New("field count mismatch")

func decode(val reflect.Value, fields *bufio.Scanner) error {
	switch val.Kind() {
	case reflect.Bool:
		return errors.New("bool not supported")
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

		for {
			// add empty element to the slice
			newElem := reflect.New(val.Type().Elem()).Elem()
			err := decode(newElem, fields)
			if err != nil {
				if err == errFieldCount {
					return nil
				}
				return err
			}
			val.Set(reflect.Append(val, newElem))
		}
	case reflect.Array:
		for i := 0; i < val.Len(); i += 1 {
			err := decode(val.Index(i), fields)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i += 1 {
			err := decode(val.Field(i), fields)
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		unwrapped := val.Elem()
		if !unwrapped.IsValid() {
			newUnwrapped := reflect.New(val.Type().Elem())
			err := decode(newUnwrapped, fields)
			if err != nil {
				return err
			}

			val.Set(newUnwrapped)
			return nil
		}

		err := decode(unwrapped, fields)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitFieldsFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.IndexByte(data, SeparatorField[0]); i >= 0 {
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
