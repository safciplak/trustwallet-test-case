package storage

import "github.com/safciplak/trustwallet/types"

type Storage interface {
	Subscribe(address string) error
	IsSubscribed(address string) (bool, error)
	SaveTransaction(address string, tx types.Transaction) error
	GetTransactions(address string) ([]types.Transaction, error)
}
