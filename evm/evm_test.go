package evm

import (
	"chaintx/chains"
	"chaintx/reporter"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChaseEvmTransfers(t *testing.T) {
	transfers := make(chan reporter.Transfer)

	go ChaseTransfers(
		chains.Ethereum,
		"0x62a7b6eb6a5d2dcaa05bf53c7272afd9da460a2c",
		20467174,
		10,
		transfers,
	)
}

func TestEvmTransfers(t *testing.T) {
	client, err := ethclient.Dial("https://eth-pokt.nodies.app")
	assert.Nil(t, err)

	transfers := transfers(
		client,
		chains.Ethereum,
		"0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		20359096,
		20359098,
	)

	assert.Equal(t, reporter.Transfer{
		Chain:        "ethereum",
		From:         "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		To:           "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		Amount:       "1",
		Txid:         "0xb831414248618c196919908bfdd05106ee7b6ab57ed8bac986374d8de7191902",
		Timestamp:    "0x884dd19c0e966eaf5b3e37e4df55a1995a243aa352e70794f0a249e7828fc274",
		TransferType: "SEND",
	}, transfers[0])
}
