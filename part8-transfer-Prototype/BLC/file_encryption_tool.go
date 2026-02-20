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
	key := deriveKey(passphrase)
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
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
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt> <input_file> <passphrase>")
		os.Exit(1)
	}

	action := os.Args[1]
	filename := os.Args[2]
	passphrase := os.Args[3]

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch action {
	case "encrypt":
		encrypted, err := encrypt(data, passphrase)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
		outputFile := filename + ".enc"
		if err := os.WriteFile(outputFile, []byte(encrypted), 0644); err != nil {
			fmt.Printf("Error writing encrypted file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted file saved as: %s\n", outputFile)

	case "decrypt":
		decrypted, err := decrypt(string(data), passphrase)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
		outputFile := filename + ".dec"
		if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
			fmt.Printf("Error writing decrypted file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted file saved as: %s\n", outputFile)

	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	saltSize   = 16
	nonceSize  = 12
	keySize    = 32
	iterations = 100000
)

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	for i := 0; i < iterations; i++ {
		hash.Write(hash.Sum(nil))
	}
	return hash.Sum(nil)[:keySize]
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(password, salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	outputData := append(salt, nonce...)
	outputData = append(outputData, ciphertext...)

	return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, password string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if len(data) < saltSize+nonceSize {
		return errors.New("file too short")
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <password>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	password := os.Args[4]

	var err error
	switch mode {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		err = errors.New("invalid mode, use 'encrypt' or 'decrypt'")
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully\n")
}
package main

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
    "path/filepath"
)

const (
    saltSize      = 16
    keyIterations = 100000
    keyLength     = 32
)

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    
    key := hash.Sum(nil)
    for i := 0; i < keyIterations-1; i++ {
        hash.Reset()
        hash.Write(key)
        key = hash.Sum(nil)
    }
    return key[:keyLength]
}

func encryptFile(inputPath, outputPath, passphrase string) error {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("failed to generate salt: %v", err)
    }

    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %v", err)
    }

    inputData, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read input file: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %v", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("failed to generate nonce: %v", err)
    }

    ciphertext := gcm.Seal(nil, nonce, inputData, nil)

    outputData := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
    outputData = append(outputData, salt...)
    outputData = append(outputData, nonce...)
    outputData = append(outputData, ciphertext...)

    if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
        return fmt.Errorf("failed to write output file: %v", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    encryptedData, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read encrypted file: %v", err)
    }

    if len(encryptedData) < saltSize {
        return errors.New("invalid encrypted file format")
    }

    salt := encryptedData[:saltSize]
    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %v", err)
    }

    nonceSize := gcm.NonceSize()
    if len(encryptedData) < saltSize+nonceSize {
        return errors.New("invalid encrypted file format")
    }

    nonce := encryptedData[saltSize : saltSize+nonceSize]
    ciphertext := encryptedData[saltSize+nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption failed: %v", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("failed to write decrypted file: %v", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Passphrase will be read from environment variable ENCRYPTION_KEY")
        os.Exit(1)
    }

    operation := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    passphrase := os.Getenv("ENCRYPTION_KEY")
    if passphrase == "" {
        fmt.Println("Error: ENCRYPTION_KEY environment variable not set")
        os.Exit(1)
    }

    var err error
    switch operation {
    case "encrypt":
        err = encryptFile(inputFile, outputFile, passphrase)
    case "decrypt":
        err = decryptFile(inputFile, outputFile, passphrase)
    default:
        fmt.Printf("Error: unknown operation '%s'\n", operation)
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", inputFile, outputFile)
}