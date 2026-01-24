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
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath, key string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, key string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return os.WriteFile(outputPath, ciphertext, 0644)
}

func generateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_encryption.go <command>")
		fmt.Println("Commands: encrypt, decrypt, generate-key")
		return
	}

	switch os.Args[1] {
	case "generate-key":
		key, err := generateKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			return
		}
		fmt.Printf("Generated key: %s\n", key)

	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run file_encryption.go encrypt <input> <output> <key>")
			return
		}
		err := encryptFile(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			return
		}
		fmt.Println("File encrypted successfully")

	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run file_encryption.go decrypt <input> <output> <key>")
			return
		}
		err := decryptFile(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			return
		}
		fmt.Println("File decrypted successfully")

	default:
		fmt.Println("Unknown command")
	}
}