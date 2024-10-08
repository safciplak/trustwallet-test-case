package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/safciplak/trustwallet/storage"
	"github.com/safciplak/trustwallet/types"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []types.Transaction
}

type PushNotificationService interface {
	Notify(transactionType string, hash, from, to, value string)
}

type LoggingPushNotificationService struct{}

func (s *LoggingPushNotificationService) Notify(transactionType string, hash, from, to, value string) {
	log.Printf("Push Notification: %s transaction - Hash: %s, From: %s, To: %s, Value: %s\n",
		transactionType, hash, from, to, value)
}

type ParserImpl struct {
	currentBlock int
	storage      storage.Storage
	rpcURL       string
	mutex        sync.Mutex

	pushNotificationService PushNotificationService
}

func NewParser(rpcURL string, storage storage.Storage) *ParserImpl {
	return &ParserImpl{
		currentBlock:            0,
		storage:                 storage,
		rpcURL:                  rpcURL,
		pushNotificationService: &LoggingPushNotificationService{},
	}
}

func (p *ParserImpl) GetCurrentBlock() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.currentBlock
}

func (p *ParserImpl) Subscribe(address string) bool {
	err := p.storage.Subscribe(address)
	return err == nil
}

func (p *ParserImpl) GetTransactions(address string) []types.Transaction {
	txs, err := p.storage.GetTransactions(address)
	if err != nil {
		fmt.Printf("Error getting transactions: %v\n", err)
		return []types.Transaction{}
	}
	return txs
}

func (p *ParserImpl) StartParsing() error {
	for {
		currentBlockNumber, err := p.GetCurrentBlockNumber()
		if err != nil {
			fmt.Printf("Error getting current block number: %v\n", err)
			return err
		}

		for i := p.currentBlock + 1; i <= currentBlockNumber; i++ {
			err := p.parseBlock(i)
			if err != nil {
				fmt.Printf("Error parsing block %d: %v\n", i, err)
				return err
			}
			p.currentBlock = i
			fmt.Printf("Parsed block %d\n", i)
		}

		time.Sleep(5 * time.Second)
	}
}

func (p *ParserImpl) GetCurrentBlockNumber() (int, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}

	respData, err := p.makeRPCRequest(data)
	if err != nil {
		return 0, err
	}

	resultHex, ok := respData["result"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	blockNumber, err := strconv.ParseInt(resultHex[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return int(blockNumber), nil
}

func (p *ParserImpl) parseBlock(blockNumber int) error {
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNumberHex, true},
		"id":      1,
	}

	respData, err := p.makeRPCRequest(data)
	if err != nil {
		return err
	}

	result, ok := respData["result"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format")
	}

	transactions, ok := result["transactions"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid transactions format")
	}

	for _, tx := range transactions {
		txMap, ok := tx.(map[string]interface{})
		if !ok {
			continue
		}

		from, _ := txMap["from"].(string)
		to, _ := txMap["to"].(string)
		hash, _ := txMap["hash"].(string)
		value, _ := txMap["value"].(string)

		fromSubscribed, err := p.storage.IsSubscribed(from)
		if err != nil {
			fmt.Printf("Error checking subscription: %v\n", err)
			continue
		}
		toSubscribed, err := p.storage.IsSubscribed(to)
		if err != nil {
			fmt.Printf("Error checking subscription: %v\n", err)
			continue
		}

		if fromSubscribed || toSubscribed {
			transaction := types.Transaction{
				Hash:  hash,
				From:  from,
				To:    to,
				Value: value,
			}

			if fromSubscribed {
				p.storage.SaveTransaction(from, transaction)
				p.pushNotificationService.Notify("Outgoing", hash, from, to, value)
			}
			if toSubscribed {
				p.storage.SaveTransaction(to, transaction)
				p.pushNotificationService.Notify("Incoming", hash, from, to, value)
			}
		}
	}

	return nil
}

func (p *ParserImpl) makeRPCRequest(data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", p.rpcURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
