package main

import (
	"fmt"

	"github.com/WhoSoup/binarymarshal"
)

type Test struct {
	A int
	B string
}

var _ binarymarshal.Marshallable = (*Test)(nil)

func (t *Test) GetMarshalOrder() []interface{} {
	return []interface{}{
		&t.A,
		&t.B,
	}
}

func main() {
	t := Test{21, "foo"}

	data, err := binarymarshal.Marshal(&t)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", data)
}
