package main

import (
	"fmt"
	"os"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
)

// Hash is a convenient fixed []byte type created at the hash length
type Hash [constants.HASH_LENGTH]byte

var _ interfaces.BinaryMarshallable = (*Hash)(nil)

// Bytes returns a copy of the internal []byte array
func (h *Hash) Bytes() (rval []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "nil hash")
			rval = constants.ZERO_HASH
		}
	}()
	return h.GetBytes()
}

// MarshalBinary returns a copy of the []byte array
func (h *Hash) MarshalBinary() (rval []byte, err error) {
	defer func(pe *error) {
		if *pe != nil {
			fmt.Fprintf(os.Stderr, "Hash.MarshalBinary err:%v", *pe)
		}
	}(&err)
	return h.Bytes(), nil
}

// UnmarshalBinaryData unmarshals the input array into the Hash
func (h *Hash) UnmarshalBinaryData(p []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()
	copy(h[:], p)
	newData = p[constants.HASH_LENGTH:]
	return
}

// UnmarshalBinary unmarshals the input array into the Hash
func (h *Hash) UnmarshalBinary(p []byte) (err error) {
	_, err = h.UnmarshalBinaryData(p)
	return
}

// GetBytes makes a copy of the hash in this hash.  Changes to the return value WILL NOT be
// reflected in the source hash.  You have to do a SetBytes to change the source
// value.
func (h *Hash) GetBytes() []byte {
	newHash := make([]byte, constants.HASH_LENGTH)
	copy(newHash, h[:])

	return newHash
}
