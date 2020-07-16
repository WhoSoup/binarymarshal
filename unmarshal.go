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
	buf := bytes.NewBuffer(data)

	if err := unmarshal(buf, o); err != nil {
		return nil
	}

	return nil
}

func unmarshal(buf *bytes.Buffer, o Marshallable) error {
	for _, field := range o.GetMarshalOrder() {
		el := reflect.ValueOf(field).Elem()
		if !el.CanSet() {
			return fmt.Errorf("unable to set field %v", el)
		}

		switch el.Kind() {
		// todo add errors for map, func?
		case reflect.Struct:
			if rec, ok := el.Interface().(Marshallable); ok {
				if err := unmarshal(buf, rec); err != nil {
					return err
				}
			} else if custom, ok := el.Interface().(Custom); ok {
				if err := custom.Decode(buf); err != nil {
					return err
				}
			} else {
				return errors.New("marshal order contains unmarshallable struct")
			}
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
		default: // see doc for binary.Write. Already works well for the most cases.
			if err := binary.Read(buf, binary.BigEndian, el.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
