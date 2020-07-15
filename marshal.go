package binarymarshal

import (
	"bytes"
	"fmt"
	"reflect"
)

func Marshal(o Marshallable) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for _, field := range o.GetMarshalOrder() {
		el := reflect.ValueOf(field).Elem()

		f := reg.getEncoder(el.Kind())
		if f == nil {
			return nil, fmt.Errorf("type %s does not have a marshaller registered", el.Type()) // todo recursive structs
		}

		fmt.Println(el, el.CanSet(), el.Type())
	}

	return buf.Bytes(), nil
}
