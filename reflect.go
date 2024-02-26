package prdemo

import (
	"errors"
	"reflect"
	"strings"
)

func (m *Message) walk(obj interface{}) error {
	return m.walkRecursive(reflect.ValueOf(obj))
}

func (m *Message) walkRecursive(val reflect.Value) error {
	r := m.r

	switch val.Kind() {
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		unwrapped := val.Elem()

		// Check if the pointer is nil
		if !unwrapped.IsValid() {
			newUnwrapped := reflect.New(val.Type().Elem())
			err := m.walkRecursive(newUnwrapped)
			if err != nil {
				return err
			}
			val.Set(newUnwrapped)
			return nil
		}

		err := m.walkRecursive(unwrapped)
		if err != nil {
			return err
		}
	case reflect.Interface:
		unwrapped := val.Elem()
		err := m.walkRecursive(unwrapped)
		if err != nil {
			return err
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i += 1 {
			// check if field has Decode method
			if d, ok := val.Interface().(Decoder); ok {
				println("has Decode method")
				err := d.Decode(m)
				if err != nil {
					return err
				}
				val.Set(reflect.ValueOf(d))
			}
			err := m.walkRecursive(val.Field(i))
			if err != nil {
				return err
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var value int64
		var err error

		switch val.Kind() {
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

		if val.CanSet() {
			val.SetInt(value)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var value uint64
		var err error

		switch val.Kind() {
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

		if val.CanSet() {
			val.SetUint(value)
		}
	case reflect.Float32:
		f, err := r.ReadFloat32()
		if err != nil {
			return err
		}

		if val.CanSet() {
			val.SetFloat(float64(f))
		}
	case reflect.Float64:
		f, err := r.ReadFloat64()
		if err != nil {
			return err
		}

		if val.CanSet() {
			val.SetFloat(f)
		}
	case reflect.Bool:
		b, err := r.ReadBool()
		if err != nil {
			return err
		}

		if val.CanSet() {
			val.SetBool(b)
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

		if val.CanSet() {
			val.SetString(builder.String())
		}
	case reflect.Slice:
		for {
			tmpV := reflect.New(val.Type().Elem()).Elem()
			err := m.walkRecursive(tmpV)
			if err != nil {
				return err
			}
			if val.CanSet() {
				val.Set(reflect.Append(val, tmpV))
			}
		}
	}

	return nil
}
