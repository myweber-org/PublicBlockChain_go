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
)

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %v", err)
	}

	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %v", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("iv generation error: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %v", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}

	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %v", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output> <passphrase>")
		fmt.Println("Example: file_encryptor encrypt secret.txt secret.enc mypassword")
		os.Exit(1)
	}

	action := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	switch action {
	case "encrypt":
		if err := encryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted successfully: %s\n", outputPath)
	case "decrypt":
		if err := decryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File decrypted successfully: %s\n", outputPath)
	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}package main

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

    sampleData := []byte("This is a secret message to be encrypted.")

    switch operation {
    case "encrypt":
        encrypted, err := encryptData(sampleData, key)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        encoded := base64.StdEncoding.EncodeToString(encrypted)
        fmt.Printf("Encrypted data (base64): %s\n", encoded)
        fmt.Printf("Encryption key (base64): %s\n", base64.StdEncoding.EncodeToString(key))

    case "decrypt":
        if len(os.Args) < 4 {
            fmt.Println("Usage for decrypt: go run file_encryptor.go decrypt <base64_ciphertext> <base64_key>")
            os.Exit(1)
        }

        ciphertext, err := base64.StdEncoding.DecodeString(os.Args[2])
        if err != nil {
            fmt.Printf("Error decoding ciphertext: %v\n", err)
            os.Exit(1)
        }

        key, err := base64.StdEncoding.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            os.Exit(1)
        }

        decrypted, err := decryptData(ciphertext, key)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted data: %s\n", decrypted)

    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
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

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt data: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output> <key>")
		os.Exit(1)
	}

	action := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	key := []byte(os.Args[4])

	if len(key) != 32 {
		fmt.Println("Key must be 32 bytes for AES-256")
		os.Exit(1)
	}

	var err error
	switch action {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, key)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, key)
	default:
		fmt.Println("Action must be 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully\n")
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

func encryptFile(inputPath, outputPath string, key []byte) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file failed: %w", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation failed: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM mode failed: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("nonce generation failed: %w", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
        return fmt.Errorf("write file failed: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file failed: %w", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation failed: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM mode failed: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption failed: %w", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write file failed: %w", err)
    }

    return nil
}

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, fmt.Errorf("key generation failed: %w", err)
    }
    return key, nil
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
        fmt.Println("Example: go run file_encryptor.go encrypt secret.txt encrypted.bin")
        os.Exit(1)
    }

    action := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]

    key, err := generateKey()
    if err != nil {
        fmt.Printf("Key generation error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Generated key: %s\n", hex.EncodeToString(key))

    switch action {
    case "encrypt":
        if err := encryptFile(inputPath, outputPath, key); err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File encrypted successfully")
    case "decrypt":
        if err := decryptFile(inputPath, outputPath, key); err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File decrypted successfully")
    default:
        fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}