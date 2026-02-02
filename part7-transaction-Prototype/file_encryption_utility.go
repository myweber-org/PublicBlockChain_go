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

    "golang.org/x/crypto/pbkdf2"
)

const (
    saltSize      = 16
    keyIterations = 10000
    keyLength     = 32
)

type EncryptionResult struct {
    Ciphertext string
    Salt       string
    IV         string
}

func deriveKey(passphrase string, salt []byte) []byte {
    return pbkdf2.Key([]byte(passphrase), salt, keyIterations, keyLength, sha256.New)
}

func encrypt(plaintext, passphrase string) (*EncryptionResult, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil, fmt.Errorf("salt generation failed: %w", err)
    }

    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("iv generation failed: %w", err)
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("cipher creation failed: %w", err)
    }

    paddedText := pkcs7Pad([]byte(plaintext), aes.BlockSize)
    ciphertext := make([]byte, len(paddedText))
    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext, paddedText)

    return &EncryptionResult{
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Salt:       base64.StdEncoding.EncodeToString(salt),
        IV:         base64.StdEncoding.EncodeToString(iv),
    }, nil
}

func decrypt(encrypted *EncryptionResult, passphrase string) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encrypted.Ciphertext)
    if err != nil {
        return "", fmt.Errorf("ciphertext decode failed: %w", err)
    }

    salt, err := base64.StdEncoding.DecodeString(encrypted.Salt)
    if err != nil {
        return "", fmt.Errorf("salt decode failed: %w", err)
    }

    iv, err := base64.StdEncoding.DecodeString(encrypted.IV)
    if err != nil {
        return "", fmt.Errorf("iv decode failed: %w", err)
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("cipher creation failed: %w", err)
    }

    if len(ciphertext)%aes.BlockSize != 0 {
        return "", errors.New("ciphertext is not a multiple of block size")
    }

    plaintext := make([]byte, len(ciphertext))
    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(plaintext, ciphertext)

    unpaddedText, err := pkcs7Unpad(plaintext)
    if err != nil {
        return "", fmt.Errorf("padding removal failed: %w", err)
    }

    return string(unpaddedText), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padText...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
    if len(data) == 0 {
        return nil, errors.New("empty data")
    }
    padding := int(data[len(data)-1])
    if padding > len(data) || padding == 0 {
        return nil, errors.New("invalid padding")
    }
    for i := 0; i < padding; i++ {
        if data[len(data)-1-i] != byte(padding) {
            return nil, errors.New("invalid padding bytes")
        }
    }
    return data[:len(data)-padding], nil
}

func main() {
    secretMessage := "Confidential data requiring secure storage"
    password := "StrongPassphrase!2024"

    fmt.Println("Original message:", secretMessage)

    encrypted, err := encrypt(secretMessage, password)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }

    fmt.Printf("\nEncryption successful\n")
    fmt.Printf("Ciphertext: %s\n", encrypted.Ciphertext[:50]+"...")
    fmt.Printf("Salt: %s\n", encrypted.Salt)
    fmt.Printf("IV: %s\n", encrypted.IV)

    decrypted, err := decrypt(encrypted, password)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }

    fmt.Printf("\nDecrypted message: %s\n", decrypted)
    fmt.Printf("Verification: %v\n", strings.EqualFold(secretMessage, decrypted))
}