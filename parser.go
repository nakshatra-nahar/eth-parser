package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Parser interface {
	GetCurrentBlock() int64
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type EthereumParser struct {
	currentBlock        int64
	subscribedAddresses map[string]bool
	addressTransactions map[string][]Transaction
	mu                  sync.Mutex
}

func NewEthereumParser() *EthereumParser {
	latestBlock, err := getLatestBlockNumber()
	if err != nil {
		log.Fatalf("Error fetching latest block number: %v", err)
	}

	const N int64 = 10000
	startBlock := latestBlock - N
	if startBlock < 0 {
		startBlock = 0
	}

	log.Printf("Starting parser from block: %d", startBlock)

	return &EthereumParser{
		currentBlock:        startBlock,
		subscribedAddresses: make(map[string]bool),
		addressTransactions: make(map[string][]Transaction),
	}
}

func (p *EthereumParser) GetCurrentBlock() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentBlock
}

func (p *EthereumParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	address = strings.ToLower(address)
	p.subscribedAddresses[address] = true
	log.Printf("Subscribed to address: %s", address)
	return true
}

func (p *EthereumParser) GetTransactions(address string) []Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()
	address = strings.ToLower(address)
	log.Printf("Fetching transactions for address: %s", address)
	return p.addressTransactions[address]
}

func (p *EthereumParser) StartParsing() {
	for {
		latestBlock, err := getLatestBlockNumber()
		if err != nil {
			log.Printf("Error getting latest block number: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("Latest block number: %d", latestBlock)

		p.mu.Lock()
		startBlock := p.currentBlock + 1
		p.mu.Unlock()

		if startBlock > latestBlock {
			log.Printf("No new blocks to process. Sleeping for 5 seconds.")
			time.Sleep(5 * time.Second)
			continue
		}

		for blockNum := startBlock; blockNum <= latestBlock; blockNum++ {
			log.Printf("Processing block number: %d", blockNum)
			block, err := getBlockByNumber(blockNum)
			if err != nil {
				log.Printf("Error getting block %d: %v", blockNum, err)
				continue
			}

			p.parseBlock(block)

			p.mu.Lock()
			p.currentBlock = blockNum
			p.mu.Unlock()

			// Add a short delay to prevent rate limiting
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (p *EthereumParser) parseBlock(block map[string]interface{}) {
	transactions, ok := block["transactions"].([]interface{})
	if !ok {
		log.Printf("No transactions found in block")
		return
	}

	timestampHex, _ := block["timestamp"].(string)
	timestamp, _ := strconv.ParseInt(strings.TrimPrefix(timestampHex, "0x"), 16, 64)

	blockNumberHex, _ := block["number"].(string)
	blockNumber, _ := strconv.ParseInt(strings.TrimPrefix(blockNumberHex, "0x"), 16, 64)

	for _, tx := range transactions {
		txMap, ok := tx.(map[string]interface{})
		if !ok {
			continue
		}

		from, _ := txMap["from"].(string)
		to, _ := txMap["to"].(string)
		hash, _ := txMap["hash"].(string)
		valueHex, _ := txMap["value"].(string)
		gasPriceHex, _ := txMap["gasPrice"].(string)
		gasUsedHex, _ := txMap["gas"].(string)

		value := hexToDecimal(valueHex)
		gasPrice := hexToDecimal(gasPriceHex)
		gasUsed := hexToDecimal(gasUsedHex)

		transactionFee := calculateTransactionFee(gasPrice, gasUsed)

		from = strings.ToLower(from)
		to = strings.ToLower(to)

		p.mu.Lock()
		_, fromSubscribed := p.subscribedAddresses[from]
		_, toSubscribed := p.subscribedAddresses[to]
		p.mu.Unlock()

		if fromSubscribed || toSubscribed {
			log.Printf("Found transaction involving subscribed address in block %d: %s", blockNumber, hash)
			transaction := Transaction{
				Hash:           hash,
				From:           from,
				To:             to,
				Value:          value,
				GasPrice:       gasPrice,
				GasUsed:        gasUsed,
				BlockNumber:    blockNumber,
				Timestamp:      timestamp,
				TransactionFee: transactionFee,
			}

			if fromSubscribed {
				p.mu.Lock()
				p.addressTransactions[from] = append(p.addressTransactions[from], transaction)
				p.mu.Unlock()
				log.Printf("Transaction added to 'from' address: %s", from)
			}
			if toSubscribed && to != from && to != "" {
				p.mu.Lock()
				p.addressTransactions[to] = append(p.addressTransactions[to], transaction)
				p.mu.Unlock()
				log.Printf("Transaction added to 'to' address: %s", to)
			}
		}
	}
}
