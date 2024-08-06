package store

import (
	"chaintx/reporter"
	"sync"
)

type TransferStore interface {
	Add(address string, transfer reporter.Transfer)
	Get(txid string) (reporter.Transfer, bool)
	ListByAddress(address string) []reporter.Transfer
}

type InMemoryTransferStore struct {
	mu                 sync.Mutex
	transfersByTxid    map[string]reporter.Transfer
	transfersByAddress map[string][]reporter.Transfer
}

func NewTransferStore() *InMemoryTransferStore {
	return &InMemoryTransferStore{
		transfersByTxid:    make(map[string]reporter.Transfer),
		transfersByAddress: make(map[string][]reporter.Transfer),
	}
}

func (s *InMemoryTransferStore) Add(address string, transfer reporter.Transfer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transfersByTxid[transfer.Txid] = transfer
	s.transfersByAddress[address] = append(s.transfersByAddress[address], transfer)
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
