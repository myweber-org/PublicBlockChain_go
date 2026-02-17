
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func calculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func verifyFileIntegrity(filePath, expectedChecksum string) (bool, error) {
	actualChecksum, err := calculateFileChecksum(filePath)
	if err != nil {
		return false, err
	}

	return actualChecksum == expectedChecksum, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: file_checksum_verifier <file_path> <expected_checksum>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	expectedChecksum := os.Args[2]

	isValid, err := verifyFileIntegrity(filePath, expectedChecksum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if isValid {
		fmt.Println("File integrity verified: checksum matches")
	} else {
		fmt.Println("WARNING: File integrity check failed - checksum mismatch")
		os.Exit(1)
	}
}