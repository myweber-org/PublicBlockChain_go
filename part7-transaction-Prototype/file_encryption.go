package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type SecureEncryptor struct {
	encryptionKey []byte
	hmacKey       []byte
}

func NewSecureEncryptor(encryptionKey, hmacKey []byte) (*SecureEncryptor, error) {
	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}
	if len(hmacKey) != 32 {
		return nil, errors.New("hmac key must be 32 bytes")
	}
	return &SecureEncryptor{
		encryptionKey: encryptionKey,
		hmacKey:       hmacKey,
	}, nil
}

func (se *SecureEncryptor) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(se.encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	mac := hmac.New(sha256.New, se.hmacKey)
	mac.Write(ciphertext)
	hmacValue := mac.Sum(nil)

	finalOutput := append(ciphertext, hmacValue...)
	return base64.StdEncoding.EncodeToString(finalOutput), nil
}

func (se *SecureEncryptor) Decrypt(encodedCiphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize+sha256.Size {
		return nil, errors.New("ciphertext too short")
	}

	ciphertext := data[:len(data)-sha256.Size]
	receivedHMAC := data[len(data)-sha256.Size:]

	mac := hmac.New(sha256.New, se.hmacKey)
	mac.Write(ciphertext)
	expectedHMAC := mac.Sum(nil)

	if !hmac.Equal(receivedHMAC, expectedHMAC) {
		return nil, errors.New("hmac verification failed")
	}

	block, err := aes.NewCipher(se.encryptionKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func generateRandomKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	encKey, _ := generateRandomKey(32)
	hmacKey, _ := generateRandomKey(32)

	encryptor, _ := NewSecureEncryptor(encKey, hmacKey)

	secretMessage := []byte("Confidential data requiring protection")
	encrypted, err := encryptor.Encrypt(secretMessage)
	if err != nil {
		fmt.Printf("Encryption error: %v\n", err)
		return
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		fmt.Printf("Decryption error: %v\n", err)
		return
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}