package binarymarshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

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
		}

		ptr := reflect.ValueOf(field) // a pointer that points to the variable that needs to be set
		el := ptr.Elem()              // the element the pointer is pointing to
		if !el.CanSet() {
			return fmt.Errorf("unable to set field %v", el)
		}

		// this is true if the order contains "&Object{...}"
		if mm, ok := field.(Marshallable); ok {
			if err := unmarshal(buf, mm.GetMarshalOrder()); err != nil {
				return err
			}
			continue
		} else if el.Kind() == reflect.Ptr {
			// check if the order contains something like "&&Object{...}" aka pointer to a pointer
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

		switch el.Kind() {
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

			el.SetString(string(data))
		case reflect.Int:
			var z int64
			if err := binary.Read(buf, binary.BigEndian, &z); err != nil {
				return err
			}
			el.SetInt(z)
		default: // see doc for binary.Read. Already works well for the most cases.
			if err := binary.Read(buf, binary.BigEndian, el.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
