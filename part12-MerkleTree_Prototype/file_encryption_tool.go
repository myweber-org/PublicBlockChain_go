package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt data: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_encryption_tool.go <command>")
		fmt.Println("Commands: encrypt, decrypt, genkey")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run file_encryption_tool.go encrypt <input> <output> <keyfile>")
			os.Exit(1)
		}
		key, err := os.ReadFile(os.Args[4])
		if err != nil {
			fmt.Printf("Error reading key: %v\n", err)
			os.Exit(1)
		}
		if err := encryptFile(os.Args[2], os.Args[3], key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")

	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run file_encryption_tool.go decrypt <input> <output> <keyfile>")
			os.Exit(1)
		}
		key, err := os.ReadFile(os.Args[4])
		if err != nil {
			fmt.Printf("Error reading key: %v\n", err)
			os.Exit(1)
		}
		if err := decryptFile(os.Args[2], os.Args[3], key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully")

	case "genkey":
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run file_encryption_tool.go genkey <keyfile>")
			os.Exit(1)
		}
		key, err := generateKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(os.Args[2], key, 0600); err != nil {
			fmt.Printf("Error writing key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Key generated successfully")

	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}