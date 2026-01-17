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
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation error: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func main() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		return
	}

	inputFile := "sensitive_data.txt"
	encryptedFile := "encrypted_data.bin"
	decryptedFile := "decrypted_data.txt"

	if err := encryptFile(inputFile, encryptedFile, key); err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	if err := decryptFile(encryptedFile, decryptedFile, key); err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")
}