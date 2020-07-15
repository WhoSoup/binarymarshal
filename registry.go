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
	reg.encoder = make(map[reflect.Type]Encoder)
	reg.decoder = make(map[reflect.Type]Decoder)
}

type registry struct {
	mtx     sync.RWMutex
	encoder map[reflect.Type]Encoder
	decoder map[reflect.Type]Decoder
}

func (r *registry) register(t reflect.Type, ef Encoder, df Decoder) {
	r.mtx.Lock()
	r.encoder[t] = ef
	r.decoder[t] = df
	r.mtx.Unlock()
}

func (r *registry) getEncoder(t reflect.Type) Encoder {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.encoder[t]
}

func (r *registry) getDecoder(t reflect.Type) Decoder {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.decoder[t]
}

func Register(t reflect.Type, ef Encoder, df Decoder) {
	reg.register(t, ef, df)
}
