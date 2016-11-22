package mom

import (
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

func TestSaveLoadAccount(t *testing.T) {
	tree := merkle.NewIAVLTree(0, nil) // in-memory
	acct := &Account{
		ID:   []byte("0123456789abcdef"),
		Name: "Jorge",
		Age:  45,
	}

	allAccts := Query{Key: AccountKey{}}
	myKey := acct.Key()
	// valid, but not in db
	otherKey := AccountKey{ID: []byte("1234567812345678")}
	// invalid key
	badKey := AccountKey{ID: []byte("foobar")}

	// check state with no data
	assert.Equal(t, 0, tree.Size())
	checkQueryCount(t, tree, allAccts, 0)
	checkLoad(t, tree, myKey, nil)
	checkLoad(t, tree, otherKey, nil)
	_, err := Load(tree, badKey)
	require.Nil(t, err)

	// save the account
	up, err := Save(tree, acct)

	assert.Equal(t, 1, tree.Size())
	require.Nil(t, err, "%+v", err)
	assert.False(t, up)
	checkQueryCount(t, tree, allAccts, 1)
	checkLoad(t, tree, myKey, acct)
	checkLoad(t, tree, otherKey, nil)
}
