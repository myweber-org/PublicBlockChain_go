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
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encrypted string, key []byte) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encrypted)
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
    if len(data) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
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
        return
    }

    switch os.Args[1] {
    case "keygen":
        key, err := generateKey()
        if err != nil {
            fmt.Printf("Key generation failed: %v\n", err)
            return
        }
        fmt.Printf("Generated key: %s\n", base64.StdEncoding.EncodeToString(key))

    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go encrypt <key_base64> <plaintext>")
            return
        }
        key, err := base64.StdEncoding.DecodeString(os.Args[2])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            return
        }
        encrypted, err := encryptData([]byte(os.Args[3]), key)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            return
        }
        fmt.Printf("Encrypted: %s\n", encrypted)

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go decrypt <key_base64> <encrypted_data>")
            return
        }
        key, err := base64.StdEncoding.DecodeString(os.Args[2])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            return
        }
        decrypted, err := decryptData(os.Args[3], key)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            return
        }
        fmt.Printf("Decrypted: %s\n", string(decrypted))

    default:
        fmt.Println("Invalid command. Use: encrypt, decrypt, or keygen")
    }
}