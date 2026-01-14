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
    "path/filepath"
)

const saltSize = 16

func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return err
    }

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return err
    }

    if _, err := outputFile.Write(salt); err != nil {
        return err
    }
    if _, err := outputFile.Write(nonce); err != nil {
        return err
    }

    plaintext, err := io.ReadAll(inputFile)
    if err != nil {
        return err
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
    _, err = outputFile.Write(ciphertext)
    return err
}

func decryptFile(inputPath, outputPath, password string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(inputFile, salt); err != nil {
        return err
    }

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(inputFile, nonce); err != nil {
        return err
    }

    ciphertext, err := io.ReadAll(inputFile)
    if err != nil {
        return err
    }

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return errors.New("decryption failed: invalid password or corrupted file")
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    _, err = outputFile.Write(plaintext)
    return err
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output> <password>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    password := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, password)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, password)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", filepath.Base(inputPath), filepath.Base(outputPath))
}