package main

import (
	"fmt"

	"github.com/FactomProject/factomd/common/primitives"
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
		&t.O,
	}
}

func main() {

	t := Test{255, "foo", &Other{666}}
	fmt.Println(t, t.O.C)

	data, err := binarymarshal.Marshal(&t)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", data)

	tt := new(Test)
	if err := binarymarshal.Unmarshal(data, tt); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n%+v\n%+v\n", t, tt, tt.O)

	data, err = binarymarshal.MarshalCustom([]interface{}{&t.A, &t.O.C})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", data)

	b := new(Test)
	b.O = new(Other)

	err = binarymarshal.UnmarshalCustom(data, []interface{}{&b.A, &b.O.C})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n%+v\n", b, b.O)

	bufff := primitives.NewBuffer(nil)
	bufff.PushInt(t.A)
	bufff.PushInt(t.O.C)
	fmt.Printf("%x\n", bufff.Bytes())
}
