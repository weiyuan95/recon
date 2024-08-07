package store

import (
	"chaintx/reporter"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Requires a running postgres instance at localhost:5432
func TestPostgresTransferStore(t *testing.T) {
	store := NewPostgresTransferStore()

	transfer := reporter.Transfer{
		Txid:         "0xb831414248618c196919908bfdd05106ee7b6ab57ed8bac986374d8de7191902",
		Timestamp:    "0x884dd19c0e966eaf5b3e37e4df55a1995a243aa352e70794f0a249e7828fc274",
		TransferType: "SEND",
		From:         "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		To:           "0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1",
		Amount:       "1",
	}
	err := store.Add("0xA5bA9D68890D0BA1C7d5c6D1AE9B2836a5c4F4f1", transfer)
	assert.Nil(t, err)
	assert.NotNil(t, store.Get("0xb831414248618c196919908bfdd05106ee7b6ab57ed8bac986374d8de7191902"))

	// Tear down
	client := getPgClient()
	//goland:noinspection SqlWithoutWhere
	_, err = client.Exec(`DELETE FROM "transfers"`)
	if err != nil {
		t.Fatal("Failed to tear down", err)
	}
}
