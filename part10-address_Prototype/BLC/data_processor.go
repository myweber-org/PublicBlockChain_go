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
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor(allowedPattern string) (*DataProcessor, error) {
	compiled, err := regexp.Compile(allowedPattern)
	if err != nil {
		return nil, err
	}
	return &DataProcessor{allowedPattern: compiled}, nil
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.allowedPattern.FindString(trimmed)
}

func (dp *DataProcessor) ValidateInput(input string) bool {
	return dp.allowedPattern.MatchString(strings.TrimSpace(input))
}

func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := dp.CleanInput(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}