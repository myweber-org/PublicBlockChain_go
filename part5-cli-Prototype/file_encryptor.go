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
    key, err := generateKey()
    if err != nil {
        fmt.Printf("Error generating key: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Generated key: %s\n", base64.StdEncoding.EncodeToString(key))

    originalText := []byte("Sensitive data requiring encryption")
    fmt.Printf("Original text: %s\n", originalText)

    encrypted, err := encryptData(originalText, key)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Encrypted (base64): %s\n", base64.StdEncoding.EncodeToString(encrypted))

    decrypted, err := decryptData(encrypted, key)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Decrypted text: %s\n", decrypted)
}