package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/primitives"
)

type RemoveServerMsg struct {
	Timestamp     *Timestamp // Message Timestamp
	ServerChainID *Hash      // ChainID of new server
	ServerType    int        // 0 = Federated, 1 = Audit

	Signature *Signature

	tyype byte
}

func (r *RemoveServerMsg) Type() byte {
	return constants.REMOVESERVER_MSG
}

func (r *RemoveServerMsg) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling Add Server Message: %v", r)
		}
		return
	}()
	newData = data
	if newData[0] != r.Type() {
		return nil, fmt.Errorf("Invalid Message type")
	}
	newData = newData[1:]

	r.Timestamp = new(Timestamp)
	newData, err = r.Timestamp.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}

	r.ServerChainID = new(Hash)
	newData, err = r.ServerChainID.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}

	r.ServerType = int(newData[0])
	newData = newData[1:]

	if len(newData) > 32 {
		r.Signature = new(Signature)
		newData, err = r.Signature.UnmarshalBinaryData(newData)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (r *RemoveServerMsg) UnmarshalBinary(data []byte) error {
	_, err := r.UnmarshalBinaryData(data)
	return err
}

func (r *RemoveServerMsg) MarshalForSignature() (rval []byte, err error) {
	defer func(pe *error) {
		if *pe != nil {
			fmt.Fprintf(os.Stderr, "RemoveServerMsg.MarshalForSignature err:%v", *pe)
		}
	}(&err)
	var buf primitives.Buffer

	binary.Write(&buf, binary.BigEndian, r.Type())

	t := r.Timestamp
	data, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = r.ServerChainID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	binary.Write(&buf, binary.BigEndian, uint8(r.ServerType))

	return buf.DeepCopyBytes(), nil
}

func (r *RemoveServerMsg) MarshalBinary() (rval []byte, err error) {
	defer func(pe *error) {
		if *pe != nil {
			fmt.Fprintf(os.Stderr, "RemoveServerMsg.MarshalBinary err:%v", *pe)
		}
	}(&err)
	var buf primitives.Buffer

	data, err := r.MarshalForSignature()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	if r.Signature != nil {
		data, err = r.Signature.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.DeepCopyBytes(), nil
}
