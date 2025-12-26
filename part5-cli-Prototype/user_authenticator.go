package main

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
    Username string `json:"username"`
    UserID   int    `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(username string, userID int, secretKey string) (string, error) {
    claims := UserClaims{
        Username: username,
        UserID:   userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "auth-service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString, secretKey string) (*UserClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

func main() {
    secretKey := "your-secret-key-here"
    
    token, err := GenerateToken("john_doe", 123, secretKey)
    if err != nil {
        fmt.Printf("Error generating token: %v\n", err)
        return
    }
    
    fmt.Printf("Generated token: %s\n", token)
    
    claims, err := ValidateToken(token, secretKey)
    if err != nil {
        fmt.Printf("Error validating token: %v\n", err)
        return
    }
    
    fmt.Printf("Token validated successfully. User: %s, ID: %d\n", claims.Username, claims.UserID)
}