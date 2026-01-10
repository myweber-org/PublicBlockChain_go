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
)

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encrypt(plaintext []byte, passphrase string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedCiphertext string, passphrase string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, err
	}

	key := deriveKey(passphrase)
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
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	secret := "Sensitive data requiring protection"
	password := "strong-password-123"

	encrypted, err := encrypt([]byte(secret), password)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := decrypt(encrypted, password)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

const saltSize = 16
const nonceSize = 12

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	result := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
	if len(ciphertext) < saltSize+nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	salt := ciphertext[:saltSize]
	nonce := ciphertext[saltSize : saltSize+nonceSize]
	actualCiphertext := ciphertext[saltSize+nonceSize:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file>")
		fmt.Println("Passphrase will be read from environment variable ENCRYPTION_PASSPHRASE")
		return
	}

	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	passphrase := os.Getenv("ENCRYPTION_PASSPHRASE")
	if passphrase == "" {
		fmt.Println("Error: ENCRYPTION_PASSPHRASE environment variable not set")
		return
	}

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	var outputData []byte
	switch operation {
	case "encrypt":
		outputData, err = encryptData(inputData, passphrase)
	case "decrypt":
		outputData, err = decryptData(inputData, passphrase)
	default:
		fmt.Println("Error: operation must be 'encrypt' or 'decrypt'")
		return
	}

	if err != nil {
		fmt.Printf("Error during %s: %v\n", operation, err)
		return
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Printf("%s operation completed successfully\n", operation)
	fmt.Printf("Input: %s (%d bytes)\n", inputFile, len(inputData))
	fmt.Printf("Output: %s (%d bytes)\n", outputFile, len(outputData))
}