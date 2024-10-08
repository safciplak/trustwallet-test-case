package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/safciplak/trustwallet/parser"
	"github.com/safciplak/trustwallet/storage"
	"github.com/safciplak/trustwallet/types"
)

func TestAPI(t *testing.T) {
	// Mock storage ve parser oluştur
	store := storage.NewMemoryStorage()
	p := parser.NewParser("mock_rpc_url", store)

	// API oluştur
	api := NewAPI(p)

	// getCurrentBlockHandler test
	t.Run("getCurrentBlockHandler", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/current_block", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.getCurrentBlockHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response map[string]int
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := response["current_block"]; !ok {
			t.Errorf("response does not contain 'current_block' key")
		}
	})

	// subscribeHandler test
	t.Run("subscribeHandler", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/subscribe?address=0x123", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.subscribeHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response map[string]string
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		if response["status"] != "subscribed" || response["address"] != "0x123" {
			t.Errorf("unexpected response: %v", response)
		}
	})

	// getTransactionsHandler test
	t.Run("getTransactionsHandler", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/transactions?address=0x123", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.getTransactionsHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response []types.Transaction
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		if response == nil {
			t.Errorf("response should not be nil")
		}
	})
}
