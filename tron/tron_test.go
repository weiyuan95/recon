package tron

import (
	"chaintx/store"
	"sync"
	"testing"
)

func TestTronTransfers(t *testing.T) {
	transfers, err := Transfers("TQCXkbg7JqBzj6q88DGJw4uBx8Gj9RftZV")
	if err != nil {
		return
	}

	// 55 as of 5th Aug, subject to change
	if len(transfers) != 55 {
		t.Fatal("Expected 55 transfers")
	}
}

func TestGoRoutineTransfers(t *testing.T) {
	transferStore := store.NewTransferStore()

	var wg sync.WaitGroup

	for _, add := range [2]string{"TQCXkbg7JqBzj6q88DGJw4uBx8Gj9RftZV", "TFgbZtBRPKdaqB4ExU9ap8m3LXvqfBDyPf"} {
		wg.Add(1)

		go func() {
			defer wg.Done()
			transfers, err := Transfers(add)
			if err != nil {
				t.Fatal("Failed to fetch transfers for TQCXkbg7JqBzj6q88DGJw4uBx8Gj9RftZV")
			}

			for _, transfer := range transfers {
				transferStore.Add(transfer.Txid, transfer)
			}
		}()
	}

	wg.Wait()
	// 113 as of 5th Aug, subject to change
	if transferStore.Length() != 113 {
		t.Fatal("Expected 113 transfers")
	}
}
