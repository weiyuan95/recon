package store

import (
	"chaintx/reporter"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryWatchedAddressStore_Add(t *testing.T) {
	store := NewWatchedAddressStore()
	expected := reporter.WatchedAddressInfo{Address: "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd"}
	store.Add("tronShasta", expected)

	fmt.Println(store.watchedAddresses)

	assert.Equal(t, 1, len(store.watchedAddresses))
	assert.Equal(t, store.watchedAddresses["tronShasta"]["TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd"], expected)
}

func TestInMemoryWatchedAddressStore_Get(t *testing.T) {
	store := NewWatchedAddressStore()
	expectedObj := reporter.WatchedAddressInfo{Address: "TPGdxSz5sFwbmrDfn7G3fjyYCJCJXPu2rd"}
	store.Add("tronShasta", expectedObj)

	actual, present := store.Get("tronShasta")

	assert.Equal(t, true, present)
	assert.Equal(t, []reporter.WatchedAddressInfo{expectedObj}, actual)
}
