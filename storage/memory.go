package storage

import (
	"strings"
	"sync"

	"github.com/safciplak/trustwallet/types"
)

type MemoryStorage struct {
	subscribers  map[string]bool
	transactions map[string][]types.Transaction
	mutex        sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		subscribers:  make(map[string]bool),
		transactions: make(map[string][]types.Transaction),
	}
}

func (ms *MemoryStorage) Subscribe(address string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.subscribers[strings.ToLower(address)] = true
	return nil
}

func (ms *MemoryStorage) IsSubscribed(address string) (bool, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	_, ok := ms.subscribers[strings.ToLower(address)]
	return ok, nil
}

func (ms *MemoryStorage) SaveTransaction(address string, tx types.Transaction) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	addr := strings.ToLower(address)
	ms.transactions[addr] = append(ms.transactions[addr], tx)
	return nil
}

func (ms *MemoryStorage) GetTransactions(address string) ([]types.Transaction, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	addr := strings.ToLower(address)
	txs, ok := ms.transactions[addr]
	if !ok {
		return []types.Transaction{}, nil
	}
	return txs, nil
}
