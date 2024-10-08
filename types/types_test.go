package types

import (
	"testing"
)

func TestTransaction(t *testing.T) {
	tx := Transaction{
		Hash:  "hashvalue",
		Value: "100",
		From:  "address1",
		To:    "address2",
	}

	if tx.Hash != "hashvalue" {
		t.Errorf("Expected Hash 'hashvalue', but got %v", tx.Hash)
	}
	if tx.Value != "100" {
		t.Errorf("Expected Value '100', but got %v", tx.Value)
	}
	if tx.From != "address1" {
		t.Errorf("Expected From 'address1', but got %v", tx.From)
	}
	if tx.To != "address2" {
		t.Errorf("Expected To 'address2', but got %v", tx.To)
	}
}
