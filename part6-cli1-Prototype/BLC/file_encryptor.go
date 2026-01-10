package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "io"
    "os"
)

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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
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
		return fmt.Errorf("gcm creation error: %w", err)
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
		return fmt.Errorf("gcm creation error: %w", err)
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

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("key generation error: %w", err)
	}
	return key, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output>")
		fmt.Println("Example: file_encryptor encrypt secret.txt secret.enc")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	key, err := generateRandomKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	switch operation {
	case "encrypt":
		if err := encryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted successfully. Key: %x\n", key)
	case "decrypt":
		fmt.Print("Enter decryption key (hex): ")
		var keyHex string
		fmt.Scanln(&keyHex)
		key = []byte(keyHex)
		if err := decryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully.")
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'.")
		os.Exit(1)
	}
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

func encryptFile(inputPath, outputPath, keyString string) error {
	key, _ := hex.DecodeString(keyString)
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256")
	}

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
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

func decryptFile(inputPath, outputPath, keyString string) error {
	key, _ := hex.DecodeString(keyString)
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256")
	}

	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
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

func generateRandomKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_encryptor <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  genkey                     Generate random encryption key")
		fmt.Println("  encrypt <input> <output>   Encrypt file")
		fmt.Println("  decrypt <input> <output>   Decrypt file")
		return
	}

	switch os.Args[1] {
	case "genkey":
		fmt.Println("Generated key:", generateRandomKey())
	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor encrypt <input> <output> <key>")
			return
		}
		err := encryptFile(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Println("Encryption failed:", err)
		} else {
			fmt.Println("File encrypted successfully")
		}
	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor decrypt <input> <output> <key>")
			return
		}
		err := decryptFile(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Println("Decryption failed:", err)
		} else {
			fmt.Println("File decrypted successfully")
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}