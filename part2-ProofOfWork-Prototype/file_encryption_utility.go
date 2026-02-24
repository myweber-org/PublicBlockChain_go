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
    nonceSize     = 12
    keyIterations = 100000
    keyLength     = 32
)

type EncryptionResult struct {
    Ciphertext string
    Salt       string
    Nonce      string
}

func deriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, keyIterations, keyLength, sha256.New)
}

func Encrypt(plaintext, password string) (*EncryptionResult, error) {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return nil, fmt.Errorf("salt generation failed: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("cipher creation failed: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("nonce generation failed: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("GCM mode initialization failed: %w", err)
    }

    ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

    return &EncryptionResult{
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Salt:       base64.StdEncoding.EncodeToString(salt),
        Nonce:      base64.StdEncoding.EncodeToString(nonce),
    }, nil
}

func Decrypt(encrypted *EncryptionResult, password string) (string, error) {
    salt, err := base64.StdEncoding.DecodeString(encrypted.Salt)
    if err != nil {
        return "", fmt.Errorf("salt decoding failed: %w", err)
    }

    nonce, err := base64.StdEncoding.DecodeString(encrypted.Nonce)
    if err != nil {
        return "", fmt.Errorf("nonce decoding failed: %w", err)
    }

    ciphertext, err := base64.StdEncoding.DecodeString(encrypted.Ciphertext)
    if err != nil {
        return "", fmt.Errorf("ciphertext decoding failed: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("cipher creation failed: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("GCM mode initialization failed: %w", err)
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", errors.New("decryption failed: invalid password or corrupted data")
    }

    return string(plaintext), nil
}

func main() {
    secretMessage := "Confidential data: API keys, tokens, and sensitive configuration"
    password := "StrongPassw0rd!2024"

    fmt.Println("Original message:", secretMessage)

    encrypted, err := Encrypt(secretMessage, password)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }

    fmt.Printf("\nEncryption Result:\n")
    fmt.Printf("Ciphertext: %s\n", encrypted.Ciphertext[:50]+"...")
    fmt.Printf("Salt: %s\n", encrypted.Salt)
    fmt.Printf("Nonce: %s\n", encrypted.Nonce)

    decrypted, err := Decrypt(encrypted, password)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }

    fmt.Printf("\nDecrypted message: %s\n", decrypted)

    if strings.Compare(secretMessage, decrypted) == 0 {
        fmt.Println("\nVerification: Encryption/decryption successful")
    } else {
        fmt.Println("\nVerification: Data mismatch detected")
    }
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)
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

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)
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
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	secretMessage := "Sensitive data requiring protection"
	passphrase := "secure-passphrase-123"

	fmt.Println("Original:", secretMessage)

	encrypted, err := encryptData([]byte(secretMessage), passphrase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encryption error: %v\n", err)
		return
	}
	fmt.Println("Encrypted (hex):", hex.EncodeToString(encrypted))

	decrypted, err := decryptData(encrypted, passphrase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Decryption error: %v\n", err)
		return
	}
	fmt.Println("Decrypted:", string(decrypted))
}