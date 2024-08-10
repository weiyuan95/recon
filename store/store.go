package store

import (
	"chaintx/chains"
	"chaintx/reporter"
	"sync"
)

type TransferStore interface {
	Add(transfer reporter.Transfer)
	Get(txid string) (reporter.Transfer, bool)
	ListByAddress(address string) []reporter.Transfer
}

// WatchedAddresses key:Address value:reporter.WatchedAddressInfo
// We use a map here since we only want to watch _unique_ addresses, there
// is no need to store duplicates
type WatchedAddresses map[string]reporter.WatchedAddressInfo

type InMemoryTransferStore struct {
	mu                 sync.Mutex
	transfersByTxid    map[string]reporter.Transfer
	transfersByAddress map[string][]reporter.Transfer
}

type InMemoryWatchedAddressStore struct {
	mu sync.Mutex
	// key:ChainName value:WatchedAddressInfo
	watchedAddresses map[chains.ChainName]WatchedAddresses
}

func NewTransferStore() *InMemoryTransferStore {
	return &InMemoryTransferStore{
		transfersByTxid:    make(map[string]reporter.Transfer),
		transfersByAddress: make(map[string][]reporter.Transfer),
	}
}

func (s *InMemoryTransferStore) Add(transfer reporter.Transfer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transfersByTxid[transfer.Txid] = transfer
	s.transfersByAddress[transfer.Address] = append(s.transfersByAddress[transfer.Address], transfer)
}

func (s *InMemoryTransferStore) Get(txid string) (reporter.Transfer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	transfer, ok := s.transfersByTxid[txid]
	return transfer, ok
}

func (s *InMemoryTransferStore) ListByAddress(address string) []reporter.Transfer {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.transfersByAddress[address]
}

func NewWatchedAddressStore() *InMemoryWatchedAddressStore {
	return &InMemoryWatchedAddressStore{
		watchedAddresses: make(map[chains.ChainName]WatchedAddresses),
	}
}

func (s *InMemoryWatchedAddressStore) Add(chainName chains.ChainName, info reporter.WatchedAddressInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// map[ChainName, map[string]reporter.WatchedAddressInfo]
	_, present := s.watchedAddresses[chainName]

	if !present {
		// Create a new map if it doesn't exist, if not we'll be trying to mutate a nil
		s.watchedAddresses[chainName] = make(WatchedAddresses)
	}

	_, present = s.watchedAddresses[chainName][info.Address]

	// Add it only if it exists
	if !present {
		s.watchedAddresses[chainName][info.Address] = info
	}
}

func (s *InMemoryWatchedAddressStore) Get(chainName chains.ChainName) ([]reporter.WatchedAddressInfo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	watchedAddressInfos, ok := s.watchedAddresses[chainName]
	// Technically inefficient and if we want to we can allocate only the required amount of space
	// but for MVP, this is good enough
	infos := make([]reporter.WatchedAddressInfo, 0)

	for _, info := range watchedAddressInfos {
		infos = append(infos, info)
	}

	return infos, ok
}

// For MVP, these global in-memory stores are good enough, we can swap this out for a persistent store later
// once the idea is validated

var LocalTransferStore = NewTransferStore()
var LocalWatchedAddressStore = NewWatchedAddressStore()
