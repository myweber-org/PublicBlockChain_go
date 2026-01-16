
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
    "strings"
)

const (
    saltSize       = 16
    nonceSize      = 12
    keyIterations  = 100000
    keyLength      = 32
)

func deriveKey(password, salt []byte) []byte {
    hash := sha256.New()
    hash.Write(password)
    hash.Write(salt)
    for i := 0; i < keyIterations-1; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keyLength]
}

func encrypt(plaintext, password []byte) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
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

func decrypt(encrypted string, password []byte) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return nil, err
    }

    if len(data) < saltSize+nonceSize {
        return nil, errors.New("invalid encrypted data")
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

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Password will be read from ENCRYPTION_PASSWORD environment variable")
        return
    }

    operation := strings.ToLower(os.Args[1])
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

    switch operation {
    case "encrypt":
        encrypted, err := encrypt(inputData, []byte(password))
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        if err := os.WriteFile(outputFile, []byte(encrypted), 0644); err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File encrypted successfully: %s\n", outputFile)

    case "decrypt":
        decrypted, err := decrypt(string(inputData), []byte(password))
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File decrypted successfully: %s\n", outputFile)

    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}