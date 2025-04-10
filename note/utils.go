package note

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

func isDigit(s string) bool {
	for _, v := range s {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func GenerateRandomHash(length int) string {
	// Implement a simple random hash generator
	// For simplicity, using a fixed string here
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range length {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate random hash: %v\n", err)
			os.Exit(1)
		}
		result[i] = charset[index.Int64()]
	}
	return string(result)
}
