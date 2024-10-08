package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/safciplak/trustwallet/parser"
	"github.com/safciplak/trustwallet/types"
)

// API provides methods to start the HTTP server.
type API struct {
	Parser *parser.ParserImpl
}

// NewAPI creates a new API instance.
func NewAPI(p *parser.ParserImpl) *API {
	return &API{
		Parser: p,
	}
}

// Start runs the HTTP server on the specified address.
func (api *API) Start(address string) {
	http.HandleFunc("/current_block", api.getCurrentBlockHandler)
	http.HandleFunc("/subscribe", api.subscribeHandler)
	http.HandleFunc("/transactions", api.getTransactionsHandler)

	fmt.Println("API server is running at", address)
	http.ListenAndServe(address, nil)
}

func (api *API) getCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	currentBlock := api.Parser.GetCurrentBlock()
	json.NewEncoder(w).Encode(map[string]int{"current_block": currentBlock})
}

func (api *API) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	api.Parser.Subscribe(address)
	json.NewEncoder(w).Encode(map[string]string{"status": "subscribed", "address": address})
}

func (api *API) getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	transactions := api.Parser.GetTransactions(address)
	if transactions == nil {
		transactions = []types.Transaction{}
	}
	json.NewEncoder(w).Encode(transactions)
}
