package mom

import (
	wutil "github.com/ethanfrey/tenderize/wire"
	"github.com/tendermint/go-merkle"
)

// Model is an abstraction over an object to be stored in MerkleDB
type Model interface {
	// Key returns the db key for this model.
	// This key may have zero values (and designed for range queries), depending on the state of the model
	// The Key should not change with any allowed transformation of the model
	Key() Key
}

// MWire is designed to bridge Models for go-wire
type MWire struct {
	Model
}

// Save attempts to save the given model in the given store
// updated is false on insert, otherwise true
// error is non-nil if key or value cannot be serialized
func Save(store merkle.Tree, model Model) (updated bool, err error) {
	key, err := KeyToBytes(model.Key())
	if err != nil {
		return false, err
	}

	data, err := ModelToBytes(model)
	if err != nil {
		return false, err
	}

	return store.Set(key, data), nil
}

// ModelToBytes converts the model into bytes to store in the db
// If there are invalid values in the model you can return an error
func ModelToBytes(model Model) ([]byte, error) {
	return wutil.ToBinary(MWire{model})
}

// ModelFromBytes sets the model contents to the passed in data
// Returns error if the data doesn't match this model
func ModelFromBytes(data []byte) (Model, error) {
	// Is there an easier way
	holder := MWire{}
	err := wutil.FromBinary(data, &holder)
	return holder.Model, err
}
