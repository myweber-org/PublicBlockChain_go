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

func main() {
    key := []byte("a-very-secret-key-32-bytes-long!")
    plaintext := "Sensitive data requiring encryption"

    encrypted, err := encrypt(key, plaintext)
    if err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        return
    }
    fmt.Printf("Encrypted: %s\n", encrypted)

    decrypted, err := decrypt(key, encrypted)
    if err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        return
    }
    fmt.Printf("Decrypted: %s\n", decrypted)
}

func encrypt(key []byte, plaintext string) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, cryptoText string) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    if len(ciphertext) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return string(ciphertext), nil
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

const keySize = 32

func generateKey() ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

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
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
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

func main() {
	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	originalText := []byte("Sensitive data requiring protection")
	fmt.Printf("Original: %s\n", originalText)

	encrypted, err := encryptData(originalText, key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Encrypted: %x\n", encrypted)

	decrypted, err := decryptData(encrypted, key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}