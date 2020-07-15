package binarymarshal

import (
	"reflect"
	"sync"
)

type Encoder func(interface{}) ([]byte, error)
type Decoder func([]byte) (interface{}, error)

var reg *registry

func init() {
	reg = new(registry)
	reg.encoder = make(map[reflect.Kind]Encoder)
	reg.decoder = make(map[reflect.Kind]Decoder)

	reg.register(reflect.Int, EncodeInt, nil)
	reg.register(reflect.String, EncodeString, nil)
}

type registry struct {
	mtx     sync.RWMutex
	encoder map[reflect.Kind]Encoder
	decoder map[reflect.Kind]Decoder
}

func (r *registry) register(t reflect.Kind, ef Encoder, df Decoder) {
	r.mtx.Lock()
	r.encoder[t] = ef
	r.decoder[t] = df
	r.mtx.Unlock()
}

func (r *registry) getEncoder(t reflect.Kind) Encoder {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.encoder[t]
}

func (r *registry) getDecoder(t reflect.Kind) Decoder {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.decoder[t]
}

func Register(t reflect.Kind, ef Encoder, df Decoder) {
	reg.register(t, ef, df)
}
