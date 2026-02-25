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

func encryptFile(inputPath, outputPath, keyHex string) error {
    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return fmt.Errorf("invalid key: %v", err)
    }
    if len(key) != 32 {
        return errors.New("key must be 32 bytes for AES-256")
    }

    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file error: %v", err)
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

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
        return fmt.Errorf("write file error: %v", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, keyHex string) error {
    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return fmt.Errorf("invalid key: %v", err)
    }
    if len(key) != 32 {
        return errors.New("key must be 32 bytes for AES-256")
    }

    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file error: %v", err)
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
        return fmt.Errorf("decryption failed: %v", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write file error: %v", err)
    }

    return nil
}

func generateRandomKey() (string, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return "", err
    }
    return hex.EncodeToString(key), nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage:")
        fmt.Println("  generatekey - generate random encryption key")
        fmt.Println("  encrypt <input> <output> <key> - encrypt file")
        fmt.Println("  decrypt <input> <output> <key> - decrypt file")
        return
    }

    switch os.Args[1] {
    case "generatekey":
        key, err := generateRandomKey()
        if err != nil {
            fmt.Printf("Error generating key: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Generated key: %s\n", key)

    case "encrypt":
        if len(os.Args) != 5 {
            fmt.Println("Usage: encrypt <input> <output> <key>")
            os.Exit(1)
        }
        if err := encryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File encrypted successfully")

    case "decrypt":
        if len(os.Args) != 5 {
            fmt.Println("Usage: decrypt <input> <output> <key>")
            os.Exit(1)
        }
        if err := decryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File decrypted successfully")

    default:
        fmt.Println("Unknown command")
        os.Exit(1)
    }
}