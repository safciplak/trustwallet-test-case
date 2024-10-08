package types

// Transaction represents an Ethereum transaction relevant to subscribed addresses.
type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
