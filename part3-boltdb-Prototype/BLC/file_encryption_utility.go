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
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 16
	nonceSize  = 12
	iterations = 100000
	keyLength  = 32
)

type EncryptionResult struct {
	Ciphertext string
	Salt       string
	Nonce      string
}

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
}

func Encrypt(plaintext, password string) (*EncryptionResult, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
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
		return "", err
	}

	nonce, err := base64.StdEncoding.DecodeString(result.Nonce)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(result.Ciphertext)
	if err != nil {
		return "", err
	}

	key := deriveKey(password, salt)

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
		return "", errors.New("decryption failed: invalid password or corrupted data")
	}

	return string(plaintext), nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <text|filepath> <password>")
		os.Exit(1)
	}

	operation := os.Args[1]
	input := os.Args[2]
	password := os.Args[3]

	var plaintext string
	if _, err := os.Stat(input); err == nil {
		data, err := os.ReadFile(input)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		plaintext = string(data)
	} else {
		plaintext = input
	}

	switch operation {
	case "encrypt":
		result, err := Encrypt(plaintext, password)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted Result:\nCiphertext: %s\nSalt: %s\nNonce: %s\n",
			result.Ciphertext, result.Salt, result.Nonce)

	case "decrypt":
		if len(os.Args) < 6 {
			fmt.Println("For decryption, provide: ciphertext salt nonce password")
			os.Exit(1)
		}
		result := &EncryptionResult{
			Ciphertext: os.Args[2],
			Salt:       os.Args[3],
			Nonce:      os.Args[4],
		}
		password := os.Args[5]

		decrypted, err := Decrypt(result, password)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted text: %s\n", decrypted)

	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
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

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)
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

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
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
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file> <passphrase>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]
	passphrase := os.Args[4]

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	var result []byte
	switch operation {
	case "encrypt":
		result, err = encryptData(inputData, passphrase)
	case "decrypt":
		result, err = decryptData(inputData, passphrase)
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error during %s operation: %v\n", operation, err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, result, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully. Output saved to %s\n", outputFile)
}