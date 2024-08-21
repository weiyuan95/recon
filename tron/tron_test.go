package tron

import (
	"github.com/stretchr/testify/assert"
	"recon/reporter"
	"recon/store"
	"sync"
	"testing"
)

func TestTrxTransfers(t *testing.T) {
	transferStore := store.NewTransferStore()

	var wg sync.WaitGroup
	for _, address := range [2]string{"TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd", "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd"} {
		wg.Add(1)

		go func() {
			defer wg.Done()
			transfers, err := TrxTransfers(address, "tronShasta")
			if err != nil {
				t.Error("Failed to fetch TRX transfers for TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd")
				return
			}

			for _, transfer := range transfers {
				transferStore.Add(transfer)
			}
		}()
	}

	wg.Wait()
	actual, present := transferStore.Get("04b33e531a668c5d71d165ae16a80ec172a86773c1c35f24956a4e08fbaef1d7")
	expected := reporter.Transfer{
		Chain:        "tronShasta",
		Address:      "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd",
		Txid:         "04b33e531a668c5d71d165ae16a80ec172a86773c1c35f24956a4e08fbaef1d7",
		Timestamp:    "1716960603000",
		TransferType: "RECEIVE",
		TokenType:    "TRX",
		From:         "TTupPaPVRejXfWMw1jSnrFqdr9mvfMFdmG",
		To:           "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd",
		Amount:       "30.000000",
	}

	assert.Equal(t, true, present)
	assert.Equal(t, expected, actual)
}

func TestTrc20Transfers(t *testing.T) {
	transferStore := store.NewTransferStore()

	var wg sync.WaitGroup
	for _, address := range [2]string{"TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd", "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd"} {
		wg.Add(1)

		go func() {
			defer wg.Done()
			transfers, err := Trc20Transfers(address, "tronShasta")
			if err != nil {
				t.Error("Failed to fetch USDT transfers for TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd")
				return
			}

			for _, transfer := range transfers {
				transferStore.Add(transfer)
			}
		}()
	}

	wg.Wait()
	actual, present := transferStore.Get("ce9865534ea5c8f3dc382a19448a574d31a2b2596d2f45c796e46eddcbdcd5b6")
	expected := reporter.Transfer{
		Chain:        "tronShasta",
		Address:      "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd",
		Txid:         "ce9865534ea5c8f3dc382a19448a574d31a2b2596d2f45c796e46eddcbdcd5b6",
		Timestamp:    "1719378372000",
		TransferType: "RECEIVE",
		TokenType:    "USDT",
		From:         "TTupPaPVRejXfWMw1jSnrFqdr9mvfMFdmG",
		To:           "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd",
		Amount:       "1.000000",
	}

	assert.Equal(t, true, present)
	assert.Equal(t, expected, actual)
}

func TestToCanonical(t *testing.T) {
	amount := "1000000"
	expected := "1.000000"

	actual, _ := toCanonical(amount, 6)

	assert.Equal(t, expected, actual)

	amount = "6808428269"
	expected = "6808.428269"

	actual, _ = toCanonical(amount, 6)
	assert.Equal(t, expected, actual)

	amount = "74485044315414"
	expected = "0.000074485044315414"

	actual, _ = toCanonical(amount, 18)
	assert.Equal(t, expected, actual)

	amount = "0"
	expected = "0.000000000000000000"

	actual, _ = toCanonical(amount, 18)
	assert.Equal(t, expected, actual)
}
