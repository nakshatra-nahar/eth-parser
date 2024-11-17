package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	parser := NewEthereumParser()

	go parser.StartParsing()

	http.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
		blockNumber := parser.GetCurrentBlock()
		log.Printf("Current block requested: %d", blockNumber)
		json.NewEncoder(w).Encode(map[string]int64{"currentBlock": blockNumber})
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			log.Printf("Invalid method for /subscribe: %s", r.Method)
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			Address string `json:"address"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data.Address == "" {
			log.Printf("Invalid request body for /subscribe: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		parser.Subscribe(data.Address)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			log.Printf("Address parameter missing in /transactions request")
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		transactions := parser.GetTransactions(address)
		if transactions == nil {
			transactions = []Transaction{}
		}
		log.Printf("Transactions returned for address %s: %d", address, len(transactions))
		json.NewEncoder(w).Encode(transactions)
	})

	http.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		txHash := r.URL.Query().Get("hash")
		if txHash == "" {
			log.Printf("Transaction hash parameter missing in /transaction request")
			http.Error(w, "Transaction hash is required", http.StatusBadRequest)
			return
		}
		transaction, err := getTransactionByHash(txHash)
		if err != nil {
			log.Printf("Error fetching transaction %s: %v", txHash, err)
			http.Error(w, fmt.Sprintf("Error fetching transaction: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(transaction)
	})

	log.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
