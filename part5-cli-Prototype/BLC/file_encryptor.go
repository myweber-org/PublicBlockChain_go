
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

func generateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }
    return b, nil
}

func encryptFile(inputPath, outputPath string, key []byte) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    iv, err := generateRandomBytes(aes.BlockSize)
    if err != nil {
        return err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    ciphertext := make([]byte, len(plaintext))
    stream.XORKeyStream(ciphertext, plaintext)

    result := append(iv, ciphertext...)
    return os.WriteFile(outputPath, result, 0644)
}

func decryptFile(inputPath, outputPath string, key []byte) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    if len(ciphertext) < aes.BlockSize {
        return errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    stream := cipher.NewCFBDecrypter(block, iv)
    plaintext := make([]byte, len(ciphertext))
    stream.XORKeyStream(plaintext, ciphertext)

    return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
        fmt.Println("Example: go run file_encryptor.go encrypt secret.txt secret.enc")
        return
    }

    operation := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]

    key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
    if len(key) != 32 {
        fmt.Println("Invalid key length. Must be 32 bytes for AES-256")
        return
    }

    switch operation {
    case "encrypt":
        err := encryptFile(inputPath, outputPath, key)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            return
        }
        fmt.Println("File encrypted successfully")
    case "decrypt":
        err := decryptFile(inputPath, outputPath, key)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            return
        }
        fmt.Println("File decrypted successfully")
    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
    }
}