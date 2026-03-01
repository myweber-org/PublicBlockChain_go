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
	"path/filepath"
)

const (
	saltSize   = 16
	nonceSize  = 12
	keySize    = 32
	iterations = 100000
)

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.Sum256([]byte(password))
	for i := 0; i < iterations; i++ {
		hash = sha256.Sum256(append(hash[:], salt...))
	}
	return hash[:keySize]
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(password, salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	outputData := append(salt, nonce...)
	outputData = append(outputData, ciphertext...)

	return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, password string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if len(ciphertext) < saltSize+nonceSize {
		return errors.New("file too short")
	}

	salt := ciphertext[:saltSize]
	nonce := ciphertext[saltSize : saltSize+nonceSize]
	ciphertext = ciphertext[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed - wrong password or corrupted file")
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <password>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	password := os.Args[4]

	var err error
	switch mode {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		fmt.Println("Mode must be 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Operation completed successfully")
}