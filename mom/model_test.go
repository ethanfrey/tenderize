package mom

import (
	"bytes"

	wutil "github.com/ethanfrey/tenderize/wire"
	"github.com/tendermint/go-wire"
)

// This file demonstrates how an application can construct there models to use this package properly

func init() {
	minAccountID = bytes.Repeat([]byte{0}, accountIDLength)
	maxAccountID = bytes.Repeat([]byte{255}, accountIDLength)
	wire.RegisterInterface(
		MWire{},
		wire.ConcreteType{O: Account{}, Byte: 0},
		// wire.ConcreteType{O: Status{}, Byte: 1},
	)
}

var (
	accountIDLength = 16
	minAccountID    []byte
	maxAccountID    []byte
)

// Account is the sample main model
type Account struct {
	ID     []byte // ID is immutable and must be 16 bytes in length
	Name   string
	Age    int32
	Status string
}

func (a *Account) Key() Key {
	return AccountKey{ID: a.ID}
}

func (a *Account) Serialize() ([]byte, error) {
	return wutil.ToBinary(MWire{a})
}

func (a *Account) Deserialize(data []byte) error {
	return wutil.FromBinary(data, MWire{a})
}

// AccountKey wraps the immutible ID
type AccountKey struct {
	ID []byte
}

func (k AccountKey) Serialize() ([]byte, error) {
	return wutil.ToBinary(k)
}

func (k AccountKey) Range() (min Key, max Key) {
	if k.ID != nil {
		return k, k
	}
	return AccountKey{ID: minAccountID}, AccountKey{ID: maxAccountID}
}

func (k AccountKey) Model() Model {
	return new(Account)
}

// // Status is the sample contained model (immutable - append only list)
// type Status struct {
// 	AccountID []byte
// 	Index     int32
// 	Message   string
// }

// type StatusKey struct {
// 	AccountID []byte
// 	Index     int32
// }
