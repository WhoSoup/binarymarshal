package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/FactomProject/factomd/common/interfaces"
)

//A structure for handling timestamps for messages
type Timestamp uint64 //in miliseconds
var _ interfaces.BinaryMarshallable = (*Timestamp)(nil)

func (t *Timestamp) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	if data == nil || len(data) < 6 {
		return nil, fmt.Errorf("Not enough data to unmarshal")
	}
	hd, data := binary.BigEndian.Uint32(data[:]), data[4:]
	ld, data := binary.BigEndian.Uint16(data[:]), data[2:]
	*t = Timestamp((uint64(hd) << 16) + uint64(ld))
	return data, nil
}

func (t *Timestamp) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (t *Timestamp) MarshalBinary() (rval []byte, err error) {
	defer func(pe *error) {
		if *pe != nil {
			fmt.Fprintf(os.Stderr, "Timestamp.MarshalBinary err:%v", *pe)
		}
	}(&err)
	var out bytes.Buffer
	hd := uint32(*t >> 16)
	ld := uint16(*t & 0xFFFF)
	binary.Write(&out, binary.BigEndian, uint32(hd))
	binary.Write(&out, binary.BigEndian, uint16(ld))
	return out.Bytes(), nil
}
