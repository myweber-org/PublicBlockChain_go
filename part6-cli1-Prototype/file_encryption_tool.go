package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	saltSize       = 16
	nonceSize      = 12
	keyIterations  = 100000
	keyLength      = 32
)

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	for i := 0; i < keyIterations-1; i++ {
		hash.Write(hash.Sum(nil))
	}
	return hash.Sum(nil)[:keyLength]
}

func encryptData(plaintext, passphrase string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	combined := append(salt, nonce...)
	combined = append(combined, ciphertext...)

	return base64.StdEncoding.EncodeToString(combined), nil
}

func decryptData(encrypted, passphrase string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < saltSize+nonceSize {
		return "", errors.New("encrypted data too short")
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func validatePassphrase(passphrase string) error {
	if len(passphrase) < 8 {
		return errors.New("passphrase must be at least 8 characters")
	}
	if !strings.ContainsAny(passphrase, "0123456789") {
		return errors.New("passphrase must contain at least one digit")
	}
	return nil
}

func main() {
	secretMessage := "Confidential data: Project X launch date 2024-12-15"
	passphrase := "SecurePass123"

	if err := validatePassphrase(passphrase); err != nil {
		fmt.Printf("Passphrase validation failed: %v\n", err)
		return
	}

	encrypted, err := encryptData(secretMessage, passphrase)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}

	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := decryptData(encrypted, passphrase)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}

	fmt.Printf("Decrypted: %s\n", decrypted)

	if decrypted == secretMessage {
		fmt.Println("Encryption/decryption successful")
	} else {
		fmt.Println("Encryption/decryption failed")
	}
}