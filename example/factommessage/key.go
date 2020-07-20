package main

import (
	"fmt"
	"os"

	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/primitives"
)

// PublicKey contains only Public part of Public/Private key pair
type PublicKey [ed25519.PublicKeySize]byte

// MarshalBinary marshals and returns a new []byte containing the PublicKey k
func (k *PublicKey) MarshalBinary() (rval []byte, err error) {
	defer func(pe *error) {
		if *pe != nil {
			fmt.Fprintf(os.Stderr, "PublicKey.MarshalBinary err:%v", *pe)
		}
	}(&err)
	var buf primitives.Buffer
	buf.Write(k[:])
	return buf.DeepCopyBytes(), nil
}

// UnmarshalBinaryData unmarshals the first section of input p of the size of a public key, and returns
// the residual data p
func (k *PublicKey) UnmarshalBinaryData(p []byte) ([]byte, error) {
	if len(p) < ed25519.PublicKeySize {
		return nil, fmt.Errorf("Invalid data passed")
	}
	copy(k[:], p)
	return p[ed25519.PublicKeySize:], nil
}

// UnmarshalBinary unmarshals the input p into the public key k
func (k *PublicKey) UnmarshalBinary(p []byte) (err error) {
	_, err = k.UnmarshalBinaryData(p)
	return
}
