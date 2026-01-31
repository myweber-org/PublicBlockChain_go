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

const keySize = 32

func generateKey() ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func encryptData(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func saveToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func readFromFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func main() {
	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		return
	}

	originalText := []byte("Confidential data requiring secure storage")
	fmt.Printf("Original: %s\n", originalText)

	encrypted, err := encryptData(originalText, key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}

	err = saveToFile("encrypted.bin", encrypted)
	if err != nil {
		fmt.Printf("Save failed: %v\n", err)
		return
	}

	loadedData, err := readFromFile("encrypted.bin")
	if err != nil {
		fmt.Printf("Read failed: %v\n", err)
		return
	}

	decrypted, err := decryptData(loadedData, key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}

	fmt.Printf("Decrypted: %s\n", decrypted)
	fmt.Printf("Verification: %v\n", string(decrypted) == string(originalText))
}