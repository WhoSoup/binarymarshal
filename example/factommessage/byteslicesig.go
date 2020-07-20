package main

import (
	"fmt"
	"os"

	"github.com/FactomProject/ed25519"
)

type ByteSliceSig [ed25519.SignatureSize]byte

// MarshalBinary marshals this ByteSliceSig into []byte array
func (bs *ByteSliceSig) MarshalBinary() (rval []byte, err error) {
	defer func(pe *error) {
		if *pe != nil {
			fmt.Fprintf(os.Stderr, "ByteSliceSig.MarshalBinary err:%v", *pe)
		}
	}(&err)
	return bs[:], nil
}

// GetFixed returns a new copy of the internal byte array
func (bs *ByteSliceSig) GetFixed() ([ed25519.SignatureSize]byte, error) {
	answer := [ed25519.SignatureSize]byte{}
	copy(answer[:], bs[:])

	return answer, nil
}

// UnmarshalBinary unmarshals the input data into the ByteSliceSig
func (bs *ByteSliceSig) UnmarshalBinary(data []byte) (err error) {
	if len(data) < ed25519.SignatureSize {
		return fmt.Errorf("Byte slice too short to unmarshal")
	}
	copy(bs[:], data[:ed25519.SignatureSize])
	return
}
