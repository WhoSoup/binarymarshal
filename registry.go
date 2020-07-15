package binarymarshal

import "reflect"

var reg *registry

func init() {
	reg = new(registry)
}

type registry struct {
	known map[reflect.Type]func()
}
