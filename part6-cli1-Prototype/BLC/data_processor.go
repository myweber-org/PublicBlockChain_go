
package main

import (
    "errors"
    "strings"
)

type UserData struct {
    Username string
    Email    string
}

func ValidateUserData(data UserData) error {
    if strings.TrimSpace(data.Username) == "" {
        return errors.New("username cannot be empty")
    }
    if !strings.Contains(data.Email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}

func TransformUsername(data UserData) UserData {
    data.Username = strings.ToLower(strings.TrimSpace(data.Username))
    return data
}

func ProcessUserInput(rawUsername, rawEmail string) (UserData, error) {
    userData := UserData{
        Username: rawUsername,
        Email:    rawEmail,
    }

    if err := ValidateUserData(userData); err != nil {
        return UserData{}, err
    }

    userData = TransformUsername(userData)
    return userData, nil
}