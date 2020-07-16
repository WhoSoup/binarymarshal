package binarymarshal

import "io"

type Custom struct {
	Encode func(io.Writer) error
	Decode func(io.Reader) error
}
