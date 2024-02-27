package prdemo

import (
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func (m *Message) walk(obj interface{}) error {
	return m.walkRecursive(reflect.ValueOf(obj))
}

func (m *Message) walkRecursive(val reflect.Value) error {
	r := m.r

	if val.Kind() != reflect.Ptr && val.Type().Implements(reflect.TypeOf((*Read)(nil)).Elem()) && val.CanInterface() {
		dec := val.Interface().(Read)
		out, err := dec.Read(m)
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(out))
		return nil
	}

	switch val.Kind() {
	case reflect.Ptr:
		return m.walkPointer(val)
	case reflect.Interface:
		unwrapped := val.Elem()
		err := m.walkRecursive(unwrapped)
		if err != nil {
			return err
		}
	case reflect.Struct:
		var flags *uint64
		for i := 0; i < val.NumField(); i += 1 {
			fieldValue := val.Field(i)

			// field tag
			fieldTag := val.Type().Field(i).Tag.Get("bin")
			if strings.Contains(fieldTag, "flag=") && flags != nil {
				f, err := strconv.Atoi(strings.Split(fieldTag, "=")[1])
				if err != nil {
					return errors.New("flag tag must be an unsigned integer")
				}
				flag := uint64(f)
				if *flags&flag == 0 {
					continue
				}
			}

			err := m.walkRecursive(fieldValue)
			if err != nil {
				return err
			}

			if fieldTag == "flags" {
				f := fieldValue.Uint()
				flags = &f
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
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return nil
				}
				return err
			}
			if val.CanSet() {
				val.Set(reflect.Append(val, tmpV))
			}
		}
	}

	return nil
}

func (m *Message) walkPointer(val reflect.Value) error {
	// To get the actual value of the original we have to call Elem()
	// At the same time this unwraps the pointer so we don't end up in
	// an infinite recursion
	unwrapped := val.Elem()

	// Check if the pointer is nil
	if !unwrapped.IsValid() {
		newUnwrapped := reflect.New(val.Type().Elem())
		if val.Type().Implements(reflect.TypeOf((*DecodeInto)(nil)).Elem()) && val.CanInterface() {
			dec := newUnwrapped.Interface().(DecodeInto)
			err := dec.Decode(m)
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(dec))
			return nil
		}

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

	return nil
}
