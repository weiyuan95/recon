package evm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/big"
	"recon/chains"
	"testing"
)

func TestGetClient(t *testing.T) {

	type testcase struct {
		chain               chains.ChainName
		expectedGenesisHash string
	}
	for _, tc := range []testcase{
		{chains.Ethereum, "0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"},
		{chains.EthereumSepolia, "0x25a5cc106eea7138acab33231d7160d69cb777ee0c2c553fcddf5138993e6dd9"},
		{chains.Polygon, "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b"},
	} {
		client, err := GetClient(tc.chain)
		assert.Nil(t, err)

		// Check that the client connects to the correct chain
		block, err := client.BlockByNumber(context.Background(), big.NewInt(0))
		assert.Nil(t, err)
		assert.Equal(t, block.Hash().Hex(), tc.expectedGenesisHash)

		// Subsequent calls should return the same client
		client2, _ := GetClient(tc.chain)
		assert.Equal(t, client, client2)
	}
}

func TestGetClient_InvalidChain(t *testing.T) {
	_, err := GetClient("invalid")
	assert.Equal(t, err.Error(), "invalid chain")
}
