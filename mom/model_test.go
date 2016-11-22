package mom

import (
	"bytes"
	"math"

	wutil "github.com/ethanfrey/tenderize/wire"
	"github.com/tendermint/go-wire"
)

// This file demonstrates how an application can construct there models to use this package properly

func init() {
	minAccountID = bytes.Repeat([]byte{0}, accountIDLength)
	maxAccountID = bytes.Repeat([]byte{255}, accountIDLength)
	wire.RegisterInterface(
		MWire{},
		wire.ConcreteType{O: &Account{}, Byte: 2},
		wire.ConcreteType{O: &Status{}, Byte: 4},
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
	// There REALLY must be an easier way
	holder := MWire{}
	err := wutil.FromBinary(data, &holder)
	if err == nil {
		*a = *(holder.Model.(*Account))
	}
	return err
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

// Status is the sample contained model (immutable - append only list)
type Status struct {
	Account AccountKey
	Index   int32
	Message string
}

func (s *Status) Key() Key {
	return StatusKey{
		Account: s.Account,
		Index:   s.Index,
	}
}

func (s *Status) Serialize() ([]byte, error) {
	return wutil.ToBinary(MWire{s})
}

func (s *Status) Deserialize(data []byte) error {
	return wutil.FromBinary(data, &MWire{s})
}

type StatusKey struct {
	Account AccountKey
	Index   int32
}

func (k StatusKey) Serialize() ([]byte, error) {
	return wutil.ToBinary(k)
}

func (k StatusKey) Range() (Key, Key) {
	a1, a2 := k.Account.Range()

	min, max := k, k
	min.Account, max.Account = a1.(AccountKey), a2.(AccountKey)

	if k.Index == 0 {
		min.Index = 1
		max.Index = math.MaxInt32
	}
	return min, max
}

func (k StatusKey) Model() Model {
	return new(Account)
}
