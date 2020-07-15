package binarymarshal

import (
	"encoding/binary"
)

// rework to use custom buffer
func EncodeInt(v interface{}) ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v.(int)))
	return buf, nil
}

func EncodeString(v interface{}) ([]byte, error) {
	return []byte(v.(string)), nil
}
