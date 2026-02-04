
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
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
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
        fmt.Println("Set environment variable ENCRYPTION_PASSPHRASE for the passphrase")
        return
    }

    passphrase := os.Getenv("ENCRYPTION_PASSPHRASE")
    if passphrase == "" {
        fmt.Println("Error: ENCRYPTION_PASSPHRASE environment variable not set")
        return
    }

    mode := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    inputData, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
        return
    }

    switch mode {
    case "encrypt":
        encrypted, err := encrypt(inputData, passphrase)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            return
        }
        err = os.WriteFile(outputFile, []byte(encrypted), 0644)
        if err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
        } else {
            fmt.Printf("File encrypted successfully to %s\n", outputFile)
        }

    case "decrypt":
        decrypted, err := decrypt(string(inputData), passphrase)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            return
        }
        err = os.WriteFile(outputFile, decrypted, 0644)
        if err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
        } else {
            fmt.Printf("File decrypted successfully to %s\n", outputFile)
        }

    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
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
)

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(passphrase, salt)

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

	outputData := append(salt, ciphertext...)

	return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	ciphertextWithSalt, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if len(ciphertextWithSalt) < 16 {
		return errors.New("file too short")
	}

	salt := ciphertextWithSalt[:16]
	ciphertext := ciphertextWithSalt[16:]

	key := deriveKey(passphrase, salt)

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
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	action := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	var err error
	switch action {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, passphrase)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, passphrase)
	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Operation completed successfully.")
}