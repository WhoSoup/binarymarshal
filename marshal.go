package binarymarshal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func Marshal(o Marshallable) ([]byte, error) {
	return MarshalCustom(o.GetMarshalOrder())
}

func MarshalCustom(order []interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := marshal(buf, order); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func marshal(buf *bytes.Buffer, order []interface{}) error {
	for _, field := range order {
		if custom, ok := field.(Custom); ok {
			if err := custom.Encode(buf); err != nil {
				return err
			}
			continue
		}

		el := reflect.ValueOf(field)
		if el.Kind() == reflect.Ptr || el.Kind() == reflect.Interface {
			el = el.Elem()
		}
		el = reflect.Indirect(el)

		if el.CanAddr() {
			if rec, ok := el.Addr().Interface().(Marshallable); ok {
				if err := marshal(buf, rec.GetMarshalOrder()); err != nil {
					return err
				}
				continue
			}
		}

		switch el.Kind() {
		// todo add errors for func, map?
		case reflect.Struct:
			// todo handle infinite loop
			//fmt.Println(el.Kind(), el.Type(), el.Interface())
			return fmt.Errorf("marshal order contains unmarshallable struct %v", el.Type())
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
