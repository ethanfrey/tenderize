package mom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-merkle"
)

func checkQueryCount(t *testing.T, tree merkle.Tree, query Query, expected int) {
	accts, err := List(tree, query)
	require.Nil(t, err)
	assert.Equal(t, expected, len(accts))
}

func checkLoad(t *testing.T, tree merkle.Tree, key Key, expected Model) {
	m, err := Load(tree, key)
	require.Nil(t, err)
	assert.EqualValues(t, expected, m)
}

func checkSave(t *testing.T, tree merkle.Tree, model Model, size int) {
	_, err := Save(tree, model)
	require.Nil(t, err, "%+v", err)
	assert.Equal(t, size, tree.Size())
}

func TestSaveLoadAccount(t *testing.T) {
	tree := merkle.NewIAVLTree(0, nil) // in-memory
	jorge := Account{
		ID:   []byte("0123456789abcdef"),
		Name: "Jorge",
		Age:  45,
	}
	maria := Account{
		ID:   []byte("42424242deadbeef"),
		Name: "Maria",
		Age:  37,
	}

	allAccts := Query{Key: AccountKey{}}
	jorgeKey := jorge.Key()
	mariaKey := maria.Key()
	// valid, but not in db
	otherKey := AccountKey{ID: []byte("1234567812345678")}
	// invalid key
	badKey := AccountKey{ID: []byte("foobar")}

	// check state with no data
	assert.Equal(t, 0, tree.Size())
	checkQueryCount(t, tree, allAccts, 0)
	checkLoad(t, tree, jorgeKey, nil)
	checkLoad(t, tree, mariaKey, nil)
	checkLoad(t, tree, otherKey, nil)
	_, err := Load(tree, badKey)
	require.Nil(t, err)

	// save the account
	up, err := Save(tree, jorge)
	assert.Equal(t, 1, tree.Size())
	require.Nil(t, err, "%+v", err)
	assert.False(t, up)

	checkQueryCount(t, tree, allAccts, 1)
	checkLoad(t, tree, jorgeKey, jorge)
	checkLoad(t, tree, mariaKey, nil)
	checkLoad(t, tree, otherKey, nil)

	// save a second account
	up, err = Save(tree, maria)
	assert.Equal(t, 2, tree.Size())
	require.Nil(t, err, "%+v", err)
	assert.False(t, up)

	checkQueryCount(t, tree, allAccts, 2)
	checkLoad(t, tree, jorgeKey, jorge)
	checkLoad(t, tree, mariaKey, maria)
	checkLoad(t, tree, otherKey, nil)

	// Jorge had a birthday
	jorge.Age = 46
	up, err = Save(tree, jorge)
	assert.Equal(t, 2, tree.Size())
	require.Nil(t, err, "%+v", err)
	assert.True(t, up)

	checkQueryCount(t, tree, allAccts, 2)
	checkLoad(t, tree, jorgeKey, jorge)
	checkLoad(t, tree, mariaKey, maria)
	checkLoad(t, tree, otherKey, nil)
}

func TestSaveLoadStatus(t *testing.T) {
	tree := merkle.NewIAVLTree(0, nil) // in-memory

	// we have two accounts for status...
	olga := Account{
		ID:   []byte("5432765454327654"),
		Name: "Olga",
		Age:  31,
	}
	vlad := Account{
		ID:   []byte("98798700asdfghjk"),
		Name: "Vladimir",
		Age:  28,
	}
	checkSave(t, tree, olga, 1)
	checkSave(t, tree, vlad, 2)

	allAccts := Query{Key: AccountKey{}}
	allStatus := Query{Key: StatusKey{Account: AccountKey{}}}
	olgaStatus := Query{Key: StatusKey{Account: olga.Key()}}
	foo, bar := allStatus.Key.Range()
	a, b := allAccts.Key.Range()
	for _, k := range []Key{allAccts.Key, a, b, allStatus.Key, olgaStatus.Key, foo, bar} {
		s, err := KeyToBytes(k)
		require.Nil(t, err)
		fmt.Println(s)
	}

	checkQueryCount(t, tree, allAccts, 2)
	checkQueryCount(t, tree, allStatus, 0)
	checkQueryCount(t, tree, olgaStatus, 0)

	os1 := Status{
		Account: olga.Key(),
		Index:   1,
		Message: "Happy",
	}
	os2 := Status{
		Account: olga.Key(),
		Index:   2,
		Message: "Sad",
	}
	vs1 := Status{
		Account: vlad.Key(),
		Index:   1,
		Message: "Say What?",
	}
	checkSave(t, tree, os1, 3)
	checkSave(t, tree, os2, 4)
	checkSave(t, tree, vs1, 5)

	checkQueryCount(t, tree, allAccts, 2)
	checkQueryCount(t, tree, allStatus, 3)
	checkQueryCount(t, tree, olgaStatus, 2)
}

func TestQueryFilters(t *testing.T) {

}
