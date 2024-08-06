package evm

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestEvmTranfers(t *testing.T) {
	client, err := ethclient.Dial("https://eth-pokt.nodies.app")
	assert.Nil(t, err)

	transfers := EvmTransfers(
		client,
		"0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		big.NewInt(20359096),
		big.NewInt(20359098),
	)

	transfer := Transfer{
		from:         "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		to:           "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		amount:       "1",
		txid:         "0xb831414248618c196919908bfdd05106ee7b6ab57ed8bac986374d8de7191902",
		timestamp:    "0x884dd19c0e966eaf5b3e37e4df55a1995a243aa352e70794f0a249e7828fc274",
		transferType: "SEND",
	}

	assert.Equal(t, transfers[0], transfer)
}
