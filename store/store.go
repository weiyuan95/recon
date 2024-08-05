package store

import (
	"chaintx/reporter"
	"sync"
)

type TransferStore struct {
	mu   sync.Mutex
	data map[string]reporter.Transfer
}

func NewTransferStore() *TransferStore {
	return &TransferStore{
		data: make(map[string]reporter.Transfer),
	}
}

func (s *TransferStore) Add(txid string, transfer reporter.Transfer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[txid] = transfer
}

func (s *TransferStore) Get(txid string) (reporter.Transfer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	transfer, ok := s.data[txid]
	return transfer, ok
}

func (s *TransferStore) Length() int {
	return len(s.data)
}
