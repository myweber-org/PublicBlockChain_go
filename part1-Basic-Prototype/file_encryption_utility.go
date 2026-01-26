
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
    "os"
)

const (
    saltSize      = 16
    nonceSize     = 12
    keyIterations = 100000
    keyLength     = 32
)

func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    for i := 0; i < keyIterations-1; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keyLength]
}

func encryptData(plaintext []byte, password string) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

    result := make([]byte, saltSize+nonceSize+len(ciphertext))
    copy(result[:saltSize], salt)
    copy(result[saltSize:saltSize+nonceSize], nonce)
    copy(result[saltSize+nonceSize:], ciphertext)

    return base64.StdEncoding.EncodeToString(result), nil
}

func decryptData(encodedCiphertext string, password string) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
    if err != nil {
        return nil, err
    }

    if len(data) < saltSize+nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    salt := data[:saltSize]
    nonce := data[saltSize : saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Password will be read from environment variable ENCRYPTION_PASSWORD")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    password := os.Getenv("ENCRYPTION_PASSWORD")
    if password == "" {
        fmt.Println("Error: ENCRYPTION_PASSWORD environment variable not set")
        os.Exit(1)
    }

    inputData, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
        os.Exit(1)
    }

    var outputData []byte
    switch mode {
    case "encrypt":
        encrypted, err := encryptData(inputData, password)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        outputData = []byte(encrypted)
    case "decrypt":
        decrypted, err := decryptData(string(inputData), password)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        outputData = decrypted
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
        fmt.Printf("Error writing output file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully. Output written to %s\n", outputFile)
}