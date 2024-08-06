package tron

import (
	"chaintx/store"
	"sync"
	"testing"
)

func TestTransfers(t *testing.T) {
	transferStore := store.NewTransferStore()

	var wg sync.WaitGroup

	for _, add := range [2]string{"TQCXkbg7JqBzj6q88DGJw4uBx8Gj9RftZV", "TFgbZtBRPKdaqB4ExU9ap8m3LXvqfBDyPf"} {
		wg.Add(1)

		go func() {
			defer wg.Done()
			transfers, err := TrxTransfers(add)
			if err != nil {
				t.Error("Failed to fetch transfers for TQCXkbg7JqBzj6q88DGJw4uBx8Gj9RftZV")
				return
			}

			for _, transfer := range transfers {
				transferStore.Add(transfer.Txid, transfer)
			}
		}()
	}

	wg.Wait()
	// 113 as of 5th Aug, subject to change
	if transferStore.Length() < 113 {
		t.Fatal("Expected at least 113 transfers")
	}
}
