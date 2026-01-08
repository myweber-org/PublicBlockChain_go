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
	saltSize   = 16
	iterations = 10000
	keyLength  = 32
)

func deriveKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, iterations, keyLength, sha256.New)
}

func EncryptText(plaintext, password string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	combined := append(salt, ciphertext...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

func DecryptText(encrypted, password string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < saltSize+aes.BlockSize {
		return "", errors.New("invalid encrypted data")
	}

	salt := data[:saltSize]
	ciphertext := data[saltSize:]

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func main() {
	password := "securePass123!"
	originalText := "Confidential data: API keys, tokens, secrets"

	encrypted, err := EncryptText(originalText, password)
	if err != nil {
		fmt.Printf("Encryption error: %v\n", err)
		return
	}

	fmt.Printf("Encrypted: %s\n", encrypted[:50]+"...")

	decrypted, err := DecryptText(encrypted, password)
	if err != nil {
		fmt.Printf("Decryption error: %v\n", err)
		return
	}

	fmt.Printf("Decrypted: %s\n", decrypted)
	fmt.Printf("Match: %v\n", strings.Compare(originalText, decrypted) == 0)
}