package prdemo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ghostiam/binstruct"
)

type unmarshal struct {
	r binstruct.Reader
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "binstruct: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "binstruct: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "binstruct: Unmarshal(nil " + e.Type.String() + ")"
}

func (u *unmarshal) Unmarshal(v interface{}) error {
	return u.unmarshal(v, nil)
}

func (u *unmarshal) unmarshal(v interface{}, parentStructValues []reflect.Value) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	structValue := rv.Elem()
	numField := structValue.NumField()

	valueType := structValue.Type()
	for i := 0; i < numField; i++ {
		fieldType := valueType.Field(i)

		fieldValue := structValue.Field(i)
		err := u.setValueToField(structValue, fieldValue, parentStructValues)
		if err != nil {
			return fmt.Errorf(`failed set value to field "%s": %w`, fieldType.Name, err)
		}
	}

	return nil
}

func (u *unmarshal) setValueToField(structValue, fieldValue reflect.Value, parentStructValues []reflect.Value) error {
	r := u.r

	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var value int64
		var err error

		switch fieldValue.Kind() {
		case reflect.Int8:
			v, e := r.ReadInt8()
			value = int64(v)
			err = e
		case reflect.Int16:
			v, e := r.ReadInt16()
			value = int64(v)
			err = e
		case reflect.Int32:
			v, e := r.ReadInt32()
			value = int64(v)
			err = e
		case reflect.Int64:
			v, e := r.ReadInt64()
			value = v
			err = e
		default: // reflect.Int:
			return errors.New("need set tag with len or use int8/int16/int32/int64")
		}

		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetInt(value)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var value uint64
		var err error

		switch fieldValue.Kind() {
		case reflect.Uint8:
			v, e := r.ReadUint8()
			value = uint64(v)
			err = e
		case reflect.Uint16:
			v, e := r.ReadUint16()
			value = uint64(v)
			err = e
		case reflect.Uint32:
			v, e := r.ReadUint32()
			value = uint64(v)
			err = e
		case reflect.Uint64:
			v, e := r.ReadUint64()
			value = v
			err = e
		default: // reflect.Uint:
			return errors.New("need set tag with len or use uint8/uint16/uint32/uint64")
		}

		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetUint(value)
		}
	case reflect.Float32:
		f, err := r.ReadFloat32()
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetFloat(float64(f))
		}
	case reflect.Float64:
		f, err := r.ReadFloat64()
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetFloat(f)
		}
	case reflect.Bool:
		b, err := r.ReadBool()
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetBool(b)
		}
	case reflect.String:
		var builder strings.Builder

		for {
			b, err := r.ReadByte()
			if err != nil {
				return err
			}

			if b == 0 {
				break
			}

			builder.WriteByte(b)
		}

		if fieldValue.CanSet() {
			fieldValue.SetString(builder.String())
		}
	case reflect.Slice:
		for {
			tmpV := reflect.New(fieldValue.Type().Elem()).Elem()
			err := u.setValueToField(structValue, tmpV, parentStructValues)
			if err != nil {
				return err
			}
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.Append(fieldValue, tmpV))
			}
		}
	// case reflect.Array:
	// 	var arrLen int64
	//
	// 	if arrLen == 0 {
	// 		arrLen = int64(fieldValue.Len())
	// 	}
	//
	// 	for i := int64(0); i < arrLen; i++ {
	// 		tmpV := reflect.New(fieldValue.Type().Elem()).Elem()
	// 		err = u.setValueToField(structValue, tmpV, fieldData.ElemFieldData, parentStructValues)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if fieldValue.CanSet() {
	// 			fieldValue.Index(int(i)).Set(tmpV)
	// 		}
	// 	}
	case reflect.Struct:
		err := u.unmarshal(fieldValue.Addr().Interface(), append(parentStructValues, structValue))
		if err != nil {
			return fmt.Errorf("unmarshal struct: %w", err)
		}
	default:
		return errors.New(`type "` + fieldValue.Kind().String() + `" not supported`)
	}

	return nil
}
