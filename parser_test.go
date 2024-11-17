package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock the Ethereum JSON-RPC responses
func mockServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}

// Test Subscribe Method
func TestSubscribe(t *testing.T) {
	parser := NewEthereumParser()
	address := "0xTestAddress"

	parser.Subscribe(address)
	if !parser.subscribedAddresses[strings.ToLower(address)] {
		t.Errorf("Address %s was not subscribed successfully", address)
	}
}

// Test GetTransactions Method
func TestGetTransactions(t *testing.T) {
	parser := NewEthereumParser()
	address := "0xTestAddress"

	// Subscribe and add a mock transaction
	parser.Subscribe(address)
	mockTransaction := Transaction{
		Hash:        "0xMockHash",
		From:        address,
		To:          "0xAnotherAddress",
		Value:       "1000000000000000000",
		BlockNumber: 1,
	}
	parser.addressTransactions[strings.ToLower(address)] = []Transaction{mockTransaction}

	transactions := parser.GetTransactions(address)
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].Hash != mockTransaction.Hash {
		t.Errorf("Expected transaction hash to be %s, got %s", mockTransaction.Hash, transactions[0].Hash)
	}
}

// Test GetCurrentBlock Method
func TestGetCurrentBlock(t *testing.T) {
	parser := NewEthereumParser()
	parser.currentBlock = 42

	if parser.GetCurrentBlock() != 42 {
		t.Errorf("Expected current block to be 42, got %d", parser.GetCurrentBlock())
	}
}

// Test hexToDecimal Utility Function
func TestHexToDecimal(t *testing.T) {
	hexStr := "0x10"
	expected := "16"

	result := hexToDecimal(hexStr)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// Test calculateTransactionFee Utility Function
func TestCalculateTransactionFee(t *testing.T) {
	gasPrice := "20000000000"        // 20 Gwei
	gasUsed := "21000"               // Gas used for a simple transfer
	expectedFee := "420000000000000" // 20 Gwei * 21000

	result := calculateTransactionFee(gasPrice, gasUsed)
	if result != expectedFee {
		t.Errorf("Expected transaction fee %s, got %s", expectedFee, result)
	}
}
