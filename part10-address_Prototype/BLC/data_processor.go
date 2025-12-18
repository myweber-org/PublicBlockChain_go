package main

import (
	"errors"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func TransformUserName(data UserData) UserData {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
	return data
}

func ProcessUserInput(rawName, rawEmail string, id int) (UserData, error) {
	user := UserData{
		ID:    id,
		Name:  rawName,
		Email: rawEmail,
	}

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	user = TransformUserName(user)
	return user, nil
}