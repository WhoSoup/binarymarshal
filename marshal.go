package binarymarshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

func Marshal(o Marshallable) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if err := marshal(buf, o); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func marshal(buf *bytes.Buffer, o Marshallable) error {
	for _, field := range o.GetMarshalOrder() {
		el := reflect.ValueOf(field).Elem()

		switch el.Kind() {
		case reflect.Func:
			if f, ok := el.Interface().(func() []byte); ok {
				if _, err := buf.Write(f()); err != nil {
					return err
				}
			} else {
				return errors.New("encountered unknown function")
			}
		case reflect.Struct:
			// todo handle infinite loop
			if rec, ok := el.Interface().(Marshallable); ok {
				if err := marshal(buf, rec); err != nil {
					return err
				}
			} else {
				return errors.New("marshal order contains unmarshallable struct")
			}
		case reflect.String:
			if err := binary.Write(buf, binary.BigEndian, []byte(el.Interface().(string))); err != nil {
				return err
			}
		case reflect.Int:
			if err := binary.Write(buf, binary.BigEndian, int64(el.Interface().(int))); err != nil {
				return err
			}
		default: // see doc for binary.Write. Already works well for the most cases.
			if err := binary.Write(buf, binary.BigEndian, el.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
