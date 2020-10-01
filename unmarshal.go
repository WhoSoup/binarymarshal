package binarymarshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var marshalType reflect.Type

func init() {
	marshalType = reflect.TypeOf((*Marshallable)(nil)).Elem()
}

func Unmarshal(data []byte, o Marshallable) error {
	return UnmarshalCustom(data, o.GetMarshalOrder())
}

func UnmarshalCustom(data []byte, order []interface{}) error {
	buf := bytes.NewBuffer(data)

	if err := unmarshal(buf, order); err != nil {
		return err
	}

	return nil

}

func unmarshal(buf *bytes.Buffer, order []interface{}) error {
	for _, field := range order {
		// custom doesn't require any reflection
		if custom, ok := field.(Custom); ok {
			if err := custom.Decode(buf); err != nil {
				return err
			}
			continue
		} else if mm, ok := field.(Marshallable); ok {
			if err := unmarshal(buf, mm.GetMarshalOrder()); err != nil {
				return err
			}
			continue
		}

		el := reflect.ValueOf(field) // a pointer that points to the variable that needs to be set

		if el.Kind() == reflect.Ptr && el.Elem().Kind() == reflect.Ptr {
			el = el.Elem()
		}

		if el.Kind() == reflect.Ptr && el.Type().Implements(marshalType) {
			// the type of el = "*Object", but we want to call new(Object)
			obj := reflect.New(el.Type().Elem()) // elem() of type = "Object"
			if marshallable, ok := obj.Interface().(Marshallable); ok {
				if err := unmarshal(buf, marshallable.GetMarshalOrder()); err != nil {
					return err
				}
				el.Set(reflect.ValueOf(marshallable))
				continue
			} else {
				return fmt.Errorf("encountered un-unmarshallable pointer of type %v", el.Type())
			}
		}

		if el.Kind() == reflect.Slice {
			if err := binary.Read(buf, binary.BigEndian, el.Interface()); err != nil {
				return err
			}
			continue
		}

		switch el.Elem().Kind() {
		case reflect.Struct:
			return errors.New("marshal order contains un-unmarshallable struct")
		case reflect.String:
			slen := new(VarInt)
			if err := slen.Unmarshal(buf); err != nil {
				return err
			}

			if buf.Len() < int(*slen) {
				return errors.New("buffer doesn't contain enough data to unmarshal string")
			}

			data := make([]byte, *slen)
			if _, err := io.ReadFull(buf, data); err != nil {
				return err
			}

			el.Elem().SetString(string(data))
		case reflect.Int:
			var z int64
			if err := binary.Read(buf, binary.BigEndian, &z); err != nil {
				return err
			}
			el.Elem().SetInt(z)
		default: // see doc for binary.Read. Already works well for the most cases.
			//fmt.Println(el, el.Kind(), el.Type())
			if err := binary.Read(buf, binary.BigEndian, el.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
