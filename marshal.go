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
		// todo add errors for func, map?
		case reflect.Struct:
			// todo handle infinite loop
			if rec, ok := el.Interface().(Marshallable); ok {
				if err := marshal(buf, rec); err != nil {
					return err
				}
			} else if custom, ok := el.Interface().(Custom); ok {
				if err := custom.Encode(buf); err != nil {
					return err
				}
			} else {
				return errors.New("marshal order contains unmarshallable struct")
			}
		case reflect.String:
			data := []byte(el.Interface().(string))
			l := VarInt(uint64(len(data)))
			if err := l.Marshal(buf); err != nil {
				return err
			}
			if err := binary.Write(buf, binary.BigEndian, data); err != nil {
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
