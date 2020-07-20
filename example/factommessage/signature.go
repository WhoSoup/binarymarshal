package main

import (
	//"encoding/hex"

	"encoding/hex"
	"fmt"
	"os"

	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/interfaces"
)

type Signature struct {
	Pub *PublicKey    `json:"pub"`
	Sig *ByteSliceSig `json:"sig"`
}

var _ interfaces.BinaryMarshallable = (*Signature)(nil)

func (e *Signature) Init() {
	if e.Pub == nil {
		e.Pub = new(PublicKey)
	}
	if e.Sig == nil {
		e.Sig = new(ByteSliceSig)
	}
}

func (s *Signature) Bytes() []byte {
	return s.Sig[:]
}

func (sig *Signature) CustomMarshalText() ([]byte, error) {
	sig.Init()
	return ([]byte)(hex.EncodeToString(sig.Pub[:]) + hex.EncodeToString(sig.Sig[:])), nil
}

func (s *Signature) MarshalBinary() (rval []byte, err error) {
	defer func(pe *error) {
		if *pe != nil {
			fmt.Fprintf(os.Stderr, "Signature.MarshalBinary err:%v", *pe)
		}
	}(&err)
	if s.Sig == nil {
		return nil, fmt.Errorf("Signature not complete")
	}
	s.Init()
	return append(s.Pub[:], s.Sig[:]...), nil
}

func (sig *Signature) UnmarshalBinaryData(data []byte) ([]byte, error) {
	if data == nil || len(data) < ed25519.SignatureSize+ed25519.PublicKeySize {
		return nil, fmt.Errorf("Not enough data to unmarshal")
	}

	sig.Sig = new(ByteSliceSig)
	var err error
	sig.Pub = new(PublicKey)
	data, err = sig.Pub.UnmarshalBinaryData(data)
	if err != nil {
		return nil, err
	}
	copy(sig.Sig[:], data[:ed25519.SignatureSize])
	data = data[ed25519.SignatureSize:]
	return data, nil
}

func (s *Signature) UnmarshalBinary(data []byte) error {
	_, err := s.UnmarshalBinaryData(data)
	return err
}
