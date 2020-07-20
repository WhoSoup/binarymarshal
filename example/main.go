package main

import (
	"fmt"

	"github.com/WhoSoup/binarymarshal"
)

type Test struct {
	A int
	B string
	O *Other
}

type Other struct {
	C int
}

func (o *Other) GetMarshalOrder() []interface{} {
	return []interface{}{&o.C}
}

var _ binarymarshal.Marshallable = (*Test)(nil)
var _ binarymarshal.Marshallable = (*Other)(nil)

func (t *Test) GetMarshalOrder() []interface{} {
	return []interface{}{
		&t.A,
		&t.B,
		t.O,
	}
}

func main() {
	t := Test{255, "foo", &Other{123}}

	data, err := binarymarshal.Marshal(&t)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", data)

	tt := new(Test)
	if err := binarymarshal.Unmarshal(data, tt); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", tt)
}
