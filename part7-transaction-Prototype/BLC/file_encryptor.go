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
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
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

func saveKeyToFile(key []byte, filename string) error {
	encodedKey := hex.EncodeToString(key)
	return os.WriteFile(filename, []byte(encodedKey), 0600)
}

func loadKeyFromFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}

func main() {
	originalText := []byte("This is a secret message that needs encryption")
	fmt.Printf("Original text: %s\n", originalText)

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	encrypted, err := encryptData(originalText, key)
	if err != nil {
		fmt.Printf("Error encrypting data: %v\n", err)
		return
	}
	fmt.Printf("Encrypted data (hex): %s\n", hex.EncodeToString(encrypted))

	decrypted, err := decryptData(encrypted, key)
	if err != nil {
		fmt.Printf("Error decrypting data: %v\n", err)
		return
	}
	fmt.Printf("Decrypted text: %s\n", decrypted)

	if string(originalText) == string(decrypted) {
		fmt.Println("Encryption and decryption successful!")
	} else {
		fmt.Println("Encryption/decryption failed!")
	}
}package main

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

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output>")
		os.Exit(1)
	}

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated key: %x\n", key)

	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	switch operation {
	case "encrypt":
		err = encryptFile(inputFile, outputFile, key)
	case "decrypt":
		err = decryptFile(inputFile, outputFile, key)
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully\n")
}
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func encryptData(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encrypted string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

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
	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	originalText := "Sensitive data requiring encryption"
	fmt.Printf("Original: %s\n", originalText)

	encrypted, err := encryptData([]byte(originalText), key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := decryptData(encrypted, key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}package main

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
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize      = 16
	nonceSize     = 12
	keyIterations = 100000
	keyLength     = 32
)

type EncryptedData struct {
	Ciphertext string `json:"ciphertext"`
	Salt       string `json:"salt"`
	Nonce      string `json:"nonce"`
}

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, keyIterations, keyLength, sha256.New)
}

func encryptData(plaintext, password string) (*EncryptedData, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	return &EncryptedData{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		Salt:       base64.StdEncoding.EncodeToString(salt),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
	}, nil
}

func decryptData(encrypted *EncryptedData, password string) (string, error) {
	salt, err := base64.StdEncoding.DecodeString(encrypted.Salt)
	if err != nil {
		return "", fmt.Errorf("invalid salt encoding: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encrypted.Nonce)
	if err != nil {
		return "", fmt.Errorf("invalid nonce encoding: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext encoding: %w", err)
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

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decryption failed: invalid password or corrupted data")
	}

	return string(plaintext), nil
}

func processFile(inputPath, outputPath, password string, encrypt bool) error {
	var data []byte
	var err error

	if inputPath == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(inputPath)
	}
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	var result string
	if encrypt {
		encrypted, err := encryptData(string(data), password)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
		result = fmt.Sprintf("SALT:%s\nNONCE:%s\nDATA:%s",
			encrypted.Salt, encrypted.Nonce, encrypted.Ciphertext)
	} else {
		lines := strings.Split(string(data), "\n")
		if len(lines) < 3 {
			return errors.New("invalid encrypted file format")
		}

		salt := strings.TrimPrefix(lines[0], "SALT:")
		nonce := strings.TrimPrefix(lines[1], "NONCE:")
		ciphertext := strings.TrimPrefix(lines[2], "DATA:")

		encrypted := &EncryptedData{
			Salt:       salt,
			Nonce:      nonce,
			Ciphertext: ciphertext,
		}

		plaintext, err := decryptData(encrypted, password)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		result = plaintext
	}

	if outputPath == "-" {
		fmt.Print(result)
	} else {
		if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input_file|-> <output_file|-> <password>")
		fmt.Println("Use '-' for stdin/stdout")
		os.Exit(1)
	}

	action := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]
	password := os.Args[4]

	encrypt := false
	switch action {
	case "encrypt":
		encrypt = true
	case "decrypt":
		encrypt = false
	default:
		fmt.Printf("Invalid action: %s. Use 'encrypt' or 'decrypt'\n", action)
		os.Exit(1)
	}

	if err := processFile(inputFile, outputFile, password, encrypt); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}