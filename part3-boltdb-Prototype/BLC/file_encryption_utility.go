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

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize      = 16
	nonceSize     = 12
	keyIterations = 100000
	keyLength     = 32
)

type EncryptionResult struct {
	Ciphertext string
	Salt       string
	Nonce      string
}

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, keyIterations, keyLength, sha256.New)
}

func Encrypt(plaintext, password string) (*EncryptionResult, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	return &EncryptionResult{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		Salt:       base64.StdEncoding.EncodeToString(salt),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
	}, nil
}

func Decrypt(result *EncryptionResult, password string) (string, error) {
	salt, err := base64.StdEncoding.DecodeString(result.Salt)
	if err != nil {
		return "", fmt.Errorf("invalid salt encoding: %w", err)
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(result.Nonce)
	if err != nil {
		return "", fmt.Errorf("invalid nonce encoding: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(result.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext encoding: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decryption failed: invalid password or corrupted data")
	}

	return string(plaintext), nil
}

func main() {
	password := "securePass123!"
	message := "Sensitive data that requires protection"

	fmt.Printf("Original message: %s\n\n", message)

	encrypted, err := Encrypt(message, password)
	if err != nil {
		fmt.Printf("Encryption error: %v\n", err)
		return
	}

	fmt.Printf("Encryption successful\n")
	fmt.Printf("Salt: %s\n", encrypted.Salt)
	fmt.Printf("Nonce: %s\n", encrypted.Nonce)
	fmt.Printf("Ciphertext: %s\n\n", encrypted.Ciphertext[:50]+"...")

	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		fmt.Printf("Decryption error: %v\n", err)
		return
	}

	fmt.Printf("Decrypted message: %s\n", decrypted)

	if strings.Compare(message, decrypted) == 0 {
		fmt.Println("\nVerification: Original and decrypted messages match")
	} else {
		fmt.Println("\nWarning: Original and decrypted messages do not match")
	}
}