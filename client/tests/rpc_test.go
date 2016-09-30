package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Make sure status is correct (we connect properly)
func TestStatus(t *testing.T) {
	c := GetClient()
	status, err := c.Status()
	if assert.Nil(t, err) {
		assert.Equal(t, GetConfig().GetString("chain_id"), status.NodeInfo.Network)
	}
}

// run most calls just to make sure no syntax errors
func TestNoErrors(t *testing.T) {
	assert := assert.New(t)
	c := GetClient()
	_, err := c.NetInfo()
	assert.Nil(err)
	_, err = c.BlockchainInfo(0, 4)
	assert.Nil(err)
	// TODO: check with a valid height
	_, err = c.Block(1000)
	assert.NotNil(err)
	// maybe this is an error???
	_, err = c.DialSeeds([]string{"one", "two"})
	assert.Nil(err)
	gen, err := c.Genesis()
	if assert.Nil(err) {
		assert.Equal(GetConfig().GetString("chain_id"), gen.Genesis.ChainID)
	}
}
