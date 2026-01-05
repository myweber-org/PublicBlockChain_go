package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

const saltSize = 16

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("salt generation failed: %w", err)
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %w", err)
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(salt); err != nil {
		return fmt.Errorf("salt write failed: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return fmt.Errorf("iv generation failed: %w", err)
	}

	if _, err := outputFile.Write(iv); err != nil {
		return fmt.Errorf("iv write failed: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := io.Copy(writer, inputFile); err != nil {
		return fmt.Errorf("encryption copy failed: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, password string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %w", err)
	}
	defer inputFile.Close()

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(inputFile, salt); err != nil {
		return fmt.Errorf("salt read failed: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inputFile, iv); err != nil {
		return fmt.Errorf("iv read failed: %w", err)
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %w", err)
	}
	defer outputFile.Close()

	reader := &cipher.StreamReader{S: stream, R: inputFile}
	if _, err := io.Copy(outputFile, reader); err != nil {
		return fmt.Errorf("decryption copy failed: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <password>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	password := os.Args[4]

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		err = errors.New("invalid operation, use 'encrypt' or 'decrypt'")
	}

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation %s completed successfully\n", operation)
}