
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_encryption.go <command>")
		fmt.Println("Commands: generate-key, encrypt <text>, decrypt <base64>")
		return
	}

	command := os.Args[1]

	switch command {
	case "generate-key":
		key, err := generateKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			return
		}
		fmt.Printf("Generated key: %s\n", base64.StdEncoding.EncodeToString(key))

	case "encrypt":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run file_encryption.go encrypt <text> <base64_key>")
			return
		}
		text := os.Args[2]
		key, err := base64.StdEncoding.DecodeString(os.Args[3])
		if err != nil {
			fmt.Printf("Invalid key: %v\n", err)
			return
		}

		ciphertext, err := encrypt([]byte(text), key)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			return
		}
		fmt.Printf("Encrypted: %s\n", base64.StdEncoding.EncodeToString(ciphertext))

	case "decrypt":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run file_encryption.go decrypt <base64_ciphertext> <base64_key>")
			return
		}
		ciphertext, err := base64.StdEncoding.DecodeString(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid ciphertext: %v\n", err)
			return
		}
		key, err := base64.StdEncoding.DecodeString(os.Args[3])
		if err != nil {
			fmt.Printf("Invalid key: %v\n", err)
			return
		}

		plaintext, err := decrypt(ciphertext, key)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			return
		}
		fmt.Printf("Decrypted: %s\n", plaintext)

	default:
		fmt.Println("Unknown command")
	}
}