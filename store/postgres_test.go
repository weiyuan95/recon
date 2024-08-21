package store

import (
	"github.com/stretchr/testify/assert"
	"recon/reporter"
	"testing"
)

// Requires a running postgres instance at localhost:5432
func TestPostgresTransferStore(t *testing.T) {
	store := NewPostgresTransferStore()
	store.Bootstrap()

	transfer := reporter.Transfer{
		Chain:        "ethereum",
		Address:      "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		Txid:         "0xb831414248618c196919908bfdd05106ee7b6ab57ed8bac986374d8de7191902",
		Timestamp:    "123",
		TransferType: "SEND",
		TokenType:    "x",
		From:         "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		To:           "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		Amount:       "1",
	}
	err := store.Add(transfer)
	assert.Nil(t, err)
	assert.NotNil(t, store.Get("0xb831414248618c196919908bfdd05106ee7b6ab57ed8bac986374d8de7191902"))

	transfer2 := reporter.Transfer{
		Chain:        "ethereum",
		Address:      "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		Txid:         "0xdef",
		Timestamp:    "123",
		TransferType: "SEND",
		TokenType:    "x",
		From:         "0xabc",
		To:           "0xdef",
		Amount:       "1",
	}
	err = store.Add(transfer2)
	assert.Nil(t, err)

	transfers := store.ListByAddress("0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1")

	assert.Equal(t, []reporter.Transfer{transfer, transfer2}, transfers)

	// Tear down
	client := getPgClient()
	//goland:noinspection SqlWithoutWhere
	_, err = client.Exec(`DELETE FROM "transfers"`)
	if err != nil {
		t.Fatal("Failed to tear down", err)
	}
}
