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
)

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encrypt(plaintext []byte, passphrase string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    result := make([]byte, len(salt)+len(ciphertext))
    copy(result[:16], salt)
    copy(result[16:], ciphertext)

    return base64.StdEncoding.EncodeToString(result), nil
}

func decrypt(encodedCiphertext string, passphrase string) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
    if err != nil {
        return nil, err
    }

    if len(data) < 16 {
        return nil, errors.New("ciphertext too short")
    }

    salt := data[:16]
    ciphertext := data[16:]

    key := deriveKey(passphrase, salt)
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
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Set environment variable ENCRYPTION_PASSPHRASE for the key")
        os.Exit(1)
    }

    passphrase := os.Getenv("ENCRYPTION_PASSPHRASE")
    if passphrase == "" {
        fmt.Println("Error: ENCRYPTION_PASSPHRASE environment variable not set")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    inputData, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
        os.Exit(1)
    }

    var outputData []byte
    switch mode {
    case "encrypt":
        encrypted, err := encrypt(inputData, passphrase)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        outputData = []byte(encrypted)
    case "decrypt":
        decrypted, err := decrypt(string(inputData), passphrase)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        outputData = decrypted
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
        fmt.Printf("Error writing output file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully. Output written to %s\n", outputFile)
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
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output>")
		return
	}

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}
	fmt.Printf("Generated key (hex): %x\n", key)

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, key)
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			return
		}
		fmt.Println("File encrypted successfully")
	case "decrypt":
		err = decryptFile(inputPath, outputPath, key)
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			return
		}
		fmt.Println("File decrypted successfully")
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
	}
}