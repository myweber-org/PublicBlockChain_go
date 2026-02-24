
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath, keyHex string) error {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("invalid key: %v", err)
	}
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	encryptedData := aesGCM.Seal(nonce, nonce, inputData, nil)

	return os.WriteFile(outputPath, encryptedData, 0644)
}

func decryptFile(inputPath, outputPath, keyHex string) error {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("invalid key: %v", err)
	}
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	encryptedData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output> <hex_key>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	keyHex := os.Args[4]

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, keyHex)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, keyHex)
	default:
		fmt.Printf("Unknown operation: %s\n", operation)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation %s completed successfully\n", operation)
}
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
		return fmt.Errorf("GCM creation error: %w", err)
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

func decryptFile(inputPath string, key []byte) ([]byte, error) {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM creation error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption error: %w", err)
	}

	return plaintext, nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("key generation error: %w", err)
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_encryptor <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  genkey                     Generate new encryption key")
		fmt.Println("  encrypt <input> <output>   Encrypt file")
		fmt.Println("  decrypt <input>            Decrypt file")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "genkey":
		key, err := generateKey()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated key: %s\n", hex.EncodeToString(key))

	case "encrypt":
		if len(os.Args) != 4 {
			fmt.Println("Usage: file_encryptor encrypt <input> <output>")
			os.Exit(1)
		}
		fmt.Print("Enter encryption key (hex): ")
		var keyHex string
		fmt.Scanln(&keyHex)
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			fmt.Printf("Invalid key format: %v\n", err)
			os.Exit(1)
		}
		if err := encryptFile(os.Args[2], os.Args[3], key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")

	case "decrypt":
		if len(os.Args) != 3 {
			fmt.Println("Usage: file_encryptor decrypt <input>")
			os.Exit(1)
		}
		fmt.Print("Enter decryption key (hex): ")
		var keyHex string
		fmt.Scanln(&keyHex)
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			fmt.Printf("Invalid key format: %v\n", err)
			os.Exit(1)
		}
		plaintext, err := decryptFile(os.Args[2], key)
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted content:\n%s\n", plaintext)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}