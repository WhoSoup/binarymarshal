package binarymarshal

import (
	"bytes"
	"io"
)

type VarInt uint64

var _ Marshallable = (*VarInt)(nil)

func (v *VarInt) GetMarshalOrder() []interface{} {
	return []interface{}{Custom{v.Marshal, v.Unmarshal}}
}

func (v *VarInt) Marshal(w io.Writer) error {
	buf := bytes.NewBuffer(nil)

	if *v == 0 {
		buf.WriteByte(0)
	}
	h := *v
	start := false

	if 0x8000000000000000&h != 0 { // Deal with the high bit set; Zero
		buf.WriteByte(0x81) // doesn't need this, only when set.
		start = true        // Going the whole 10 byte path!
	}

	for i := 0; i < 9; i++ {
		b := byte(h >> 56) // Get the top 7 bits
		if b != 0 || start {
			start = true
			if i != 8 {
				b = b | 0x80
			} else {
				b = b & 0x7F
			}
			buf.WriteByte(b)
		}
		h = h << 7
	}

	_, err := w.Write(buf.Bytes())
	return err
}
func (v *VarInt) Unmarshal(r io.Reader) error {
	var x uint64
	buf := make([]byte, 1)

	for {
		if _, err := r.Read(buf); err != nil {
			return err
		}
		x = x << 7
		x += uint64(buf[0]) & 0x7F
		if buf[0] < 0x80 {
			break
		}
	}
	*v = VarInt(x)
	return nil
}
