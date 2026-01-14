
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
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
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

func encryptData(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func saveKeyToFile(key []byte, filename string) error {
	encodedKey := hex.EncodeToString(key)
	return os.WriteFile(filename, []byte(encodedKey), 0600)
}

func loadKeyFromFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	if len(key) != keySize {
		return nil, errors.New("invalid key size")
	}

	return key, nil
}

func main() {
	originalText := []byte("This is a secret message that needs encryption.")
	fmt.Printf("Original text: %s\n", originalText)

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	err = saveKeyToFile(key, "encryption.key")
	if err != nil {
		fmt.Printf("Error saving key: %v\n", err)
		return
	}
	fmt.Println("Encryption key saved to 'encryption.key'")

	encrypted, err := encryptData(originalText, key)
	if err != nil {
		fmt.Printf("Error encrypting data: %v\n", err)
		return
	}
	fmt.Printf("Encrypted data (hex): %s\n", hex.EncodeToString(encrypted))

	loadedKey, err := loadKeyFromFile("encryption.key")
	if err != nil {
		fmt.Printf("Error loading key: %v\n", err)
		return
	}

	decrypted, err := decryptData(encrypted, loadedKey)
	if err != nil {
		fmt.Printf("Error decrypting data: %v\n", err)
		return
	}
	fmt.Printf("Decrypted text: %s\n", decrypted)
}