package mom

import (
	"bytes"
	"math"

	"github.com/tendermint/go-wire"
)

// This file demonstrates how an application can construct there models to use this package properly

func init() {
	minAccountID = bytes.Repeat([]byte{0}, accountIDLength)
	maxAccountID = bytes.Repeat([]byte{255}, accountIDLength)
	wire.RegisterInterface(
		MWire{},
		wire.ConcreteType{O: Account{}, Byte: 1},
		wire.ConcreteType{O: Status{}, Byte: 2},
	)
	wire.RegisterInterface(
		MKey{},
		wire.ConcreteType{O: AccountKey{}, Byte: 1},
		wire.ConcreteType{O: StatusKey{}, Byte: 2},
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

// AccountKey wraps the immutible ID
type AccountKey struct {
	ID []byte
}

func (a Account) Key() Key {
	return AccountKey{ID: a.ID}
}

func (k AccountKey) Range() (min Key, max Key) {
	if len(k.ID) == accountIDLength {
		return k, k
	}
	// TODO: if len > 0 but < 16, then use the prefix and fill the rest with 0 or 255 for min, max
	return AccountKey{ID: minAccountID}, AccountKey{ID: maxAccountID}
}

// Status is the sample contained model (immutable - append only list)
type Status struct {
	Account Key
	Index   int32
	Message string
}

type StatusKey struct {
	Account Key
	Index   int32
}

func (s Status) Key() Key {
	return StatusKey{
		Account: s.Account,
		Index:   s.Index,
	}
}

func (k StatusKey) Range() (Key, Key) {
	// TODO: make this a bit cleaner?
	min, max := k, k
	min.Account, max.Account = k.Account.Range()

	if k.Index == 0 {
		min.Index = 1
		max.Index = math.MaxInt32
	}
	return min, max
}
