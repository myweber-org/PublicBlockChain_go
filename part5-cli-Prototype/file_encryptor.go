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

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, keyHex string) error {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("invalid key: %v", err)
	}
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256")
	}

	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %v", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output> <key_hex>")
		fmt.Println("Example key: 6368616e676520746869732070617373776f726420746f206120736563726574")
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
	key := []byte("32-byte-long-key-here-for-aes-256!!")
	
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
		return
	}

	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	switch operation {
	case "encrypt":
		if err := encryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
		} else {
			fmt.Println("File encrypted successfully")
		}
	case "decrypt":
		if err := decryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
		} else {
			fmt.Println("File decrypted successfully")
		}
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
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

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 16)
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

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return append(salt, ciphertext...), nil
}

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
	if len(ciphertext) < 16 {
		return nil, errors.New("ciphertext too short")
	}

	salt := ciphertext[:16]
	ciphertext = ciphertext[16:]

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
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input_file> <output_file>")
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

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	var outputData []byte
	switch operation {
	case "encrypt":
		outputData, err = encryptData(inputData, passphrase)
	case "decrypt":
		outputData, err = decryptData(inputData, passphrase)
	default:
		fmt.Printf("Invalid operation: %s\n", operation)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error during %s: %v\n", operation, err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation %s completed successfully\n", operation)
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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt|keygen>")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "keygen":
        key, err := generateKey()
        if err != nil {
            fmt.Printf("Error generating key: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Generated key: %s\n", base64.StdEncoding.EncodeToString(key))

    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go encrypt <input_file> <base64_key>")
            os.Exit(1)
        }

        key, err := base64.StdEncoding.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            os.Exit(1)
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            os.Exit(1)
        }

        encrypted, err := encryptData(data, key)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }

        err = os.WriteFile(os.Args[2]+".enc", encrypted, 0644)
        if err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted file saved as: %s.enc\n", os.Args[2])

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go decrypt <encrypted_file> <base64_key>")
            os.Exit(1)
        }

        key, err := base64.StdEncoding.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            os.Exit(1)
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            os.Exit(1)
        }

        decrypted, err := decryptData(data, key)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }

        outputFile := os.Args[2]
        if len(outputFile) > 4 && outputFile[len(outputFile)-4:] == ".enc" {
            outputFile = outputFile[:len(outputFile)-4]
        }
        outputFile = outputFile + ".dec"

        err = os.WriteFile(outputFile, decrypted, 0644)
        if err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted file saved as: %s\n", outputFile)

    default:
        fmt.Println("Invalid command. Use: encrypt, decrypt, or keygen")
        os.Exit(1)
    }
}package main

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

func generateRandomKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt>")
        os.Exit(1)
    }

    operation := os.Args[1]
    key, err := generateRandomKey()
    if err != nil {
        fmt.Printf("Error generating key: %v\n", err)
        os.Exit(1)
    }

    sampleData := []byte("This is a secret message that needs protection.")

    switch operation {
    case "encrypt":
        encrypted, err := encryptData(sampleData, key)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted (hex): %s\n", hex.EncodeToString(encrypted))
        fmt.Printf("Key (hex): %s\n", hex.EncodeToString(key))

    case "decrypt":
        if len(os.Args) < 4 {
            fmt.Println("Usage for decrypt: go run file_encryptor.go decrypt <hex_ciphertext> <hex_key>")
            os.Exit(1)
        }
        ciphertext, err := hex.DecodeString(os.Args[2])
        if err != nil {
            fmt.Printf("Invalid ciphertext: %v\n", err)
            os.Exit(1)
        }
        key, err := hex.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            os.Exit(1)
        }
        decrypted, err := decryptData(ciphertext, key)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted: %s\n", decrypted)

    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}