package store

import (
	"chaintx/reporter"
	"sync"
)

type TransferStore struct {
	mu                 sync.Mutex
	transfersByTxid    map[string]reporter.Transfer
	transfersByAddress map[string][]reporter.Transfer
}

func NewTransferStore() *TransferStore {
	return &TransferStore{
		transfersByTxid:    make(map[string]reporter.Transfer),
		transfersByAddress: make(map[string][]reporter.Transfer),
	}
}

func (s *TransferStore) Add(address string, transfer reporter.Transfer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transfersByTxid[transfer.Txid] = transfer
	s.transfersByAddress[address] = append(s.transfersByAddress[address], transfer)
}

func (s *TransferStore) Get(txid string) (reporter.Transfer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	transfer, ok := s.transfersByTxid[txid]
	return transfer, ok
}

func (s *TransferStore) ListByAddress(address string) []reporter.Transfer {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.transfersByAddress[address]
}
