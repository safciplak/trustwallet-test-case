package storage

import (
	"testing"

	"github.com/safciplak/trustwallet/types"
)

type MockStorage struct {
	subscriptions map[string]bool
	transactions  map[string][]types.Transaction
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		subscriptions: make(map[string]bool),
		transactions:  make(map[string][]types.Transaction),
	}
}

func (ms *MockStorage) Subscribe(address string) error {
	ms.subscriptions[address] = true
	return nil
}

func (ms *MockStorage) IsSubscribed(address string) (bool, error) {
	subscribed, exists := ms.subscriptions[address]
	if !exists {
		return false, nil
	}
	return subscribed, nil
}

func (ms *MockStorage) SaveTransaction(address string, tx types.Transaction) error {
	ms.transactions[address] = append(ms.transactions[address], tx)
	return nil
}

func (ms *MockStorage) GetTransactions(address string) ([]types.Transaction, error) {
	return ms.transactions[address], nil
}

func TestStorageInterface(t *testing.T) {
	storage := NewMockStorage()
	address := "testAddress"
	tx := types.Transaction{
		// Transaction details
	}

	// Test Subscribe function
	err := storage.Subscribe(address)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Test IsSubscribed function
	subscribed, err := storage.IsSubscribed(address)
	if err != nil {
		t.Fatalf("IsSubscribed failed: %v", err)
	}
	if !subscribed {
		t.Errorf("Address should have been subscribed")
	}

	// Test SaveTransaction function
	err = storage.SaveTransaction(address, tx)
	if err != nil {
		t.Fatalf("SaveTransaction failed: %v", err)
	}

	// Test GetTransactions function
	transactions, err := storage.GetTransactions(address)
	if err != nil {
		t.Fatalf("GetTransactions failed: %v", err)
	}
	if len(transactions) != 1 || transactions[0] != tx {
		t.Errorf("Transactions do not match expected values")
	}
}
