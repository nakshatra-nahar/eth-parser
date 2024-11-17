package main

type Transaction struct {
	Hash           string `json:"hash"`
	From           string `json:"from"`
	To             string `json:"to"`
	Value          string `json:"value"`
	GasPrice       string `json:"gasPrice"`
	GasUsed        string `json:"gasUsed"`
	BlockNumber    int64  `json:"blockNumber"`
	Timestamp      int64  `json:"timestamp"`
	TransactionFee string `json:"transactionFee"`
}
