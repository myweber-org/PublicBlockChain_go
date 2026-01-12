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
        fmt.Println("Example: go run file_encryptor.go encrypt secret.txt secret.enc")
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

    keyHex := hex.EncodeToString(key)
    fmt.Printf("Generated key (hex): %s\n", keyHex)
    fmt.Println("Save this key for decryption!")

    switch action {
    case "encrypt":
        if err := encryptFile(inputPath, outputPath, key); err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File encrypted successfully")
    case "decrypt":
        fmt.Print("Enter decryption key (hex): ")
        var keyHexInput string
        fmt.Scanln(&keyHexInput)

        keyInput, err := hex.DecodeString(keyHexInput)
        if err != nil {
            fmt.Printf("Invalid key format: %v\n", err)
            os.Exit(1)
        }

        if err := decryptFile(inputPath, outputPath, keyInput); err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File decrypted successfully")
    default:
        fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}