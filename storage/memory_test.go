package storage

import (
	"testing"

	"github.com/safciplak/trustwallet/types"
)

func TestMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()

	address := "0xTestAddress"

	// Test Subscribe
	err := ms.Subscribe(address)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Test IsSubscribed
	subscribed, err := ms.IsSubscribed(address)
	if err != nil {
		t.Fatalf("IsSubscribed failed: %v", err)
	}
	if !subscribed {
		t.Fatalf("Address should be subscribed")
	}

	// Test SaveTransaction
	tx := types.Transaction{
		Hash:  "0xHash",
		From:  "0xFrom",
		To:    "0xTo",
		Value: "0xValue",
	}

	err = ms.SaveTransaction(address, tx)
	if err != nil {
		t.Fatalf("SaveTransaction failed: %v", err)
	}

	// Test GetTransactions
	txs, err := ms.GetTransactions(address)
	if err != nil {
		t.Fatalf("GetTransactions failed: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(txs))
	}
	if txs[0].Hash != tx.Hash {
		t.Fatalf("Transaction hash mismatch")
	}
}
