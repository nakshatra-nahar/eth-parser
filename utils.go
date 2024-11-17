package main

import (
	"math/big"
	"strings"
)

func hexToDecimal(hexStr string) string {
	if hexStr == "" {
		return "0"
	}
	value, ok := new(big.Int).SetString(strings.TrimPrefix(hexStr, "0x"), 16)
	if !ok {
		return "0"
	}
	return value.String()
}

func calculateTransactionFee(gasPriceStr, gasUsedStr string) string {
	gasPrice, ok := new(big.Int).SetString(gasPriceStr, 10)
	if !ok {
		return "0"
	}
	gasUsed, ok := new(big.Int).SetString(gasUsedStr, 10)
	if !ok {
		return "0"
	}
	fee := new(big.Int).Mul(gasPrice, gasUsed)
	return fee.String()
}

func weiToEther(weiStr string) string {
	wei, ok := new(big.Int).SetString(weiStr, 10)
	if !ok {
		return "0"
	}
	etherValue := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return etherValue.Text('f', 18)
}
