package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath, key string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, key string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return os.WriteFile(outputPath, ciphertext, 0644)
}

func generateKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_encryptor <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  genkey - generate a random encryption key")
		fmt.Println("  encrypt <input> <output> <key> - encrypt a file")
		fmt.Println("  decrypt <input> <output> <key> - decrypt a file")
		return
	}

	switch os.Args[1] {
	case "genkey":
		fmt.Println("Generated key:", generateKey())
	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor encrypt <input> <output> <key>")
			return
		}
		if err := encryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			fmt.Println("Encryption failed:", err)
		} else {
			fmt.Println("File encrypted successfully")
		}
	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor decrypt <input> <output> <key>")
			return
		}
		if err := decryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			fmt.Println("Decryption failed:", err)
		} else {
			fmt.Println("File decrypted successfully")
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}