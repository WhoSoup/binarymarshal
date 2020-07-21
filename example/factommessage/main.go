package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/WhoSoup/binarymarshal"
)

func (m *RemoveServerMsg) GetSignatureMarshalOrder() []interface{} {
	return []interface{}{
		&m.tyype,
		&m.Timestamp,
		&m.ServerChainID,
		&m.ServerType,
	}
}

func (m *RemoveServerMsg) GetMarshalOrder() []interface{} {
	return append(m.GetSignatureMarshalOrder(), &m.Signature)
}

func (ts *Timestamp) GetMarshalOrder() []interface{} {
	return []interface{}{
		binarymarshal.Custom{
			Encode: func(w io.Writer) error {
				hd := uint32(*ts >> 16)
				ld := uint16(*ts & 0xFFFF)
				binary.Write(w, binary.BigEndian, uint32(hd))
				binary.Write(w, binary.BigEndian, uint16(ld))
				return nil
			},
			Decode: func(r io.Reader) error {
				var hd uint32
				var ld uint16
				if err := binary.Read(r, binary.BigEndian, &hd); err != nil {
					return err
				}
				if err := binary.Read(r, binary.BigEndian, &ld); err != nil {
					return err
				}
				*ts = Timestamp((uint64(hd) << 16) + uint64(ld))
				return nil
			},
		},
	}
}

func (s *Signature) GetMarshalOrder() []interface{} {
	return []interface{}{
		&s.Pub,
		&s.Sig,
	}
}

func (pk *PublicKey) GetMarshalOrder() []interface{} {
	return []interface{}{
		pk[:],
	}
}

func (bss *ByteSliceSig) GetMarshalOrder() []interface{} {
	return []interface{}{
		bss[:],
	}
}

func (h *Hash) GetMarshalOrder() []interface{} {
	return []interface{}{
		h[:],
	}
}

func main() {
	msg := createMessage()

	// factom's method of binary marshal
	data1, err := msg.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("old: %x\n", data1)

	// new method of marshal
	data2, err := binarymarshal.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new: %x\n", data2)

	// new method of unmarshal
	msg2 := new(RemoveServerMsg)
	if err := binarymarshal.Unmarshal(data2, msg2); err != nil {
		log.Fatal(err)
	}

	fmt.Println("equal?",
		msg.tyype == msg2.tyype,
		bytes.Equal(msg.ServerChainID[:], msg2.ServerChainID[:]),
		msg.ServerType == msg2.ServerType,
		bytes.Equal(msg.Signature.Pub[:], msg2.Signature.Pub[:]),
		bytes.Equal(msg.Signature.Sig[:], msg2.Signature.Sig[:]),
	)

}

func createMessage() *RemoveServerMsg {
	msg := new(RemoveServerMsg)
	// 2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae
	chain, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	msg.ServerChainID = new(Hash)
	copy(msg.ServerChainID[:], chain)

	msg.ServerType = 1

	sig := new(Signature)
	sig.Pub = new(PublicKey)
	for i := 0; i < len(sig.Pub); i++ {
		sig.Pub[i] = 0xbb
	}
	sig.Sig = new(ByteSliceSig)
	for i := 0; i < len(sig.Sig); i++ {
		sig.Sig[i] = 0xcc
	}
	msg.Signature = sig

	msg.Timestamp = new(Timestamp)
	*msg.Timestamp = 0x1111111111111111
	msg.tyype = constants.REMOVESERVER_MSG

	return msg
}
