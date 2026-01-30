package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const saltSize = 16

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	if _, err := output.Write(salt); err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	if _, err := output.Write(iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: output}

	if _, err := io.Copy(writer, input); err != nil {
		return err
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(input, salt); err != nil {
		return err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(input, iv); err != nil {
		return err
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	stream := cipher.NewCTR(block, iv)
	reader := &cipher.StreamReader{S: stream, R: input}

	if _, err := io.Copy(output, reader); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	var err error
	switch mode {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, passphrase)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, passphrase)
	default:
		err = errors.New("invalid mode, use 'encrypt' or 'decrypt'")
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully: %s -> %s\n", filepath.Base(inputPath), filepath.Base(outputPath))
}