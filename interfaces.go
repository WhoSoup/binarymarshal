package binarymarshal

type Marshallable interface {
	GetMarshalOrder() []*interface{}
}
