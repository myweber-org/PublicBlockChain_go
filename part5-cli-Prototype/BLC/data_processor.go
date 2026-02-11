
package main

import "fmt"

func FilterAndDoublePositiveInts(nums []int) []int {
    var result []int
    for _, num := range nums {
        if num > 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{-5, 2, 0, 8, -1, 3}
    output := FilterAndDoublePositiveInts(input)
    fmt.Printf("Input: %v\n", input)
    fmt.Printf("Output: %v\n", output)
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// ValidateJSONString checks if the provided string is valid JSON.
func ValidateJSONString(input string) bool {
    var js interface{}
    return json.Unmarshal([]byte(input), &js) == nil
}

// PrettyPrintJSON takes a JSON string and prints it in a formatted, indented way.
func PrettyPrintJSON(jsonStr string) error {
    if !ValidateJSONString(jsonStr) {
        return fmt.Errorf("invalid JSON string")
    }

    var data interface{}
    if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
        return err
    }

    prettyJSON, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    fmt.Println(string(prettyJSON))
    return nil
}

// ExtractJSONField attempts to extract a top-level string field from a JSON object.
func ExtractJSONField(jsonStr, fieldName string) (string, error) {
    if !ValidateJSONString(jsonStr) {
        return "", fmt.Errorf("invalid JSON string")
    }

    var result map[string]interface{}
    if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
        return "", err
    }

    if value, exists := result[fieldName]; exists {
        if str, ok := value.(string); ok {
            return str, nil
        }
        return "", fmt.Errorf("field '%s' is not a string", fieldName)
    }
    return "", fmt.Errorf("field '%s' not found", fieldName)
}

func main() {
    sampleJSON := `{"name":"test","active":true,"count":42}`

    fmt.Println("Valid JSON?", ValidateJSONString(sampleJSON))

    fmt.Println("\nPretty printed JSON:")
    PrettyPrintJSON(sampleJSON)

    fmt.Println("\nExtracting 'name' field:")
    if name, err := ExtractJSONField(sampleJSON, "name"); err == nil {
        fmt.Println("Found name:", strings.ToUpper(name))
    }
}