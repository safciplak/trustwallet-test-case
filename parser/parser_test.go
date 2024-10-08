package parser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/safciplak/trustwallet/storage"
	"github.com/safciplak/trustwallet/types"
)

// MockStorage implements the Storage interface for testing purposes.
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

func (m *MockStorage) Subscribe(address string) error {
	m.subscriptions[address] = true
	return nil
}

func (m *MockStorage) IsSubscribed(address string) (bool, error) {
	return m.subscriptions[address], nil
}

func (m *MockStorage) SaveTransaction(address string, tx types.Transaction) error {
	m.transactions[address] = append(m.transactions[address], tx)
	return nil
}

func (m *MockStorage) GetTransactions(address string) ([]types.Transaction, error) {
	return m.transactions[address], nil
}

func TestGetCurrentBlockNumber(t *testing.T) {
	// Mock the RPC server response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "0x64", // Hexadecimal for 100
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	p := NewParser(server.URL, storage.NewMemoryStorage())

	blockNum, err := p.GetCurrentBlockNumber()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if blockNum != 100 {
		t.Errorf("expected block number 100, got %v", blockNum)
	}
}

func TestParseBlock(t *testing.T) {
	// Mock the RPC server response for eth_getBlockByNumber
	blockNumber := 100
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)
		method := reqBody["method"]
		var response map[string]interface{}
		if method == "eth_getBlockByNumber" {
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result": map[string]interface{}{
					"transactions": []interface{}{
						map[string]interface{}{
							"from":  "0xabc",
							"to":    "0xdef",
							"hash":  "0x123",
							"value": "0x1",
						},
					},
				},
			}
		} else if method == "eth_blockNumber" {
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result":  "0x" + strconv.FormatInt(int64(blockNumber), 16),
			}
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	mockStorage := NewMockStorage()
	mockStorage.Subscribe("0xabc")
	p := NewParser(server.URL, mockStorage)

	err := p.parseBlock(blockNumber)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	txs, err := mockStorage.GetTransactions("0xabc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %v", len(txs))
	}
	if txs[0].Hash != "0x123" {
		t.Errorf("expected transaction hash '0x123', got '%v'", txs[0].Hash)
	}
}
