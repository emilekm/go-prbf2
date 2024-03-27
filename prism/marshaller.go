package prism

import (
	"reflect"
	"strconv"
)

func Marshal(m any) Message {
	var subject Subject
	switch m.(type) {
	case Login1Request:
		subject = SubjectLogin1
	case Login2Request:
		subject = SubjectLogin2
	case RACommand:
		subject = SubjectRACommand
	}

	msg := Message{
		Subject: subject,
	}

	val := reflect.ValueOf(m)
	for i := 0; i < val.NumField(); i += 1 {
		field := val.Field(i)
		fieldType := field.Type()

		switch fieldType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValueInt64 := field.Int()
			msg.AddField(StringToBytes(strconv.Itoa(int(fieldValueInt64))))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValueUint64 := field.Uint()
			msg.AddField(StringToBytes(strconv.Itoa(int(fieldValueUint64))))
		case reflect.String:
			fieldValueString := field.String()
			msg.AddField([]byte(fieldValueString))
		case reflect.Slice:
			if fieldType.Elem().Kind() == reflect.Uint8 {
				msg.AddField(field.Bytes())
			}
		}
	}

	return msg
}
