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
    "strings"
)

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encrypt(plaintext, passphrase string) (string, error) {
    salt := make([]byte, 16)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
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

    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    ciphertext = append(salt, ciphertext...)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedCiphertext, passphrase string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
    if err != nil {
        return "", err
    }

    if len(data) < 16 {
        return "", errors.New("ciphertext too short")
    }

    salt := data[:16]
    ciphertext := data[16:]
    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func main() {
    secretMessage := "Confidential data: API keys, tokens, and credentials"
    password := "securePass123!"

    fmt.Println("Original:", secretMessage)

    encrypted, err := encrypt(secretMessage, password)
    if err != nil {
        fmt.Println("Encryption error:", err)
        return
    }
    fmt.Println("Encrypted:", encrypted[:50]+"...")

    decrypted, err := decrypt(encrypted, password)
    if err != nil {
        fmt.Println("Decryption error:", err)
        return
    }
    fmt.Println("Decrypted:", decrypted)

    testWrongPass, _ := decrypt(encrypted, "wrongPassword")
    fmt.Println("Wrong password test:", strings.Contains(testWrongPass, "Confidential"))
}