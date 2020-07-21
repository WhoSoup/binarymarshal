package main

import (
	"testing"

	"github.com/WhoSoup/binarymarshal"
)

func BenchmarkMarshal(b *testing.B) {
	msg := createMessage()

	b.Run("factomd method", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := msg.MarshalBinary(); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("new method", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := binarymarshal.Marshal(msg); err != nil {
				b.Error(err)
			}
		}
	})

	data, err := msg.MarshalBinary()
	if err != nil {
		b.Error(err)
	}

	b.Run("factomd method unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := msg.UnmarshalBinary(data); err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("new method unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := binarymarshal.Unmarshal(data, msg); err != nil {
				b.Error(err)
			}
		}
	})

}
