package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getLatestBlockNumber() (int64, error) {
	url := "https://ethereum-rpc.publicnode.com"
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}

	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var response struct {
		Result string `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}

	blockNumber, err := strconv.ParseInt(strings.TrimPrefix(response.Result, "0x"), 16, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func getBlockByNumber(blockNumber int64) (map[string]interface{}, error) {
	url := "https://ethereum-rpc.publicnode.com"
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNumberHex, true},
		"id":      1,
	}

	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Result map[string]interface{} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if response.Result == nil {
		return nil, fmt.Errorf("Block %d not found", blockNumber)
	}

	return response.Result, nil
}

func getTransactionByHash(txHash string) (*Transaction, error) {
	url := "https://ethereum-rpc.publicnode.com"
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByHash",
		"params":  []interface{}{txHash},
		"id":      1,
	}

	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Result map[string]interface{} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	txMap := response.Result
	if txMap == nil {
		return nil, fmt.Errorf("Transaction not found")
	}

	from, _ := txMap["from"].(string)
	to, _ := txMap["to"].(string)
	hash, _ := txMap["hash"].(string)
	valueHex, _ := txMap["value"].(string)
	gasPriceHex, _ := txMap["gasPrice"].(string)
	gasUsedHex, _ := txMap["gas"].(string)
	blockNumberHex, _ := txMap["blockNumber"].(string)

	blockNumber, err := strconv.ParseInt(strings.TrimPrefix(blockNumberHex, "0x"), 16, 64)
	if err != nil {
		return nil, err
	}
	block, err := getBlockByNumber(blockNumber)
	if err != nil {
		return nil, err
	}
	timestampHex, _ := block["timestamp"].(string)
	timestamp, _ := strconv.ParseInt(strings.TrimPrefix(timestampHex, "0x"), 16, 64)

	value := hexToDecimal(valueHex)
	gasPrice := hexToDecimal(gasPriceHex)
	gasUsed := hexToDecimal(gasUsedHex)

	transactionFee := calculateTransactionFee(gasPrice, gasUsed)

	transaction := &Transaction{
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

	return transaction, nil
}
