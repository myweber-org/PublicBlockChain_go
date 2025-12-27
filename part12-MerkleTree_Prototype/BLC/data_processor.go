
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// DataPayload represents a generic structure for incoming JSON data.
type DataPayload struct {
    ID    int             `json:"id"`
    Value string          `json:"value"`
    Tags  []string        `json:"tags"`
    Meta  json.RawMessage `json:"meta"`
}

// ValidatePayload checks the basic integrity of a DataPayload.
func ValidatePayload(payload *DataPayload) error {
    if payload.ID <= 0 {
        return fmt.Errorf("invalid ID: must be positive integer")
    }
    if strings.TrimSpace(payload.Value) == "" {
        return fmt.Errorf("value cannot be empty or whitespace")
    }
    return nil
}

// ParseJSONData attempts to unmarshal raw JSON bytes into a DataPayload and validate it.
func ParseJSONData(rawData []byte) (*DataPayload, error) {
    var payload DataPayload
    if err := json.Unmarshal(rawData, &payload); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    if err := ValidatePayload(&payload); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    return &payload, nil
}

func main() {
    // Example usage
    jsonStr := `{"id": 42, "value": "sample data", "tags": ["go", "json"], "meta": {"version": 1}}`
    data := []byte(jsonStr)

    result, err := ParseJSONData(data)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Parsed payload: ID=%d, Value=%s, Tags=%v\n", result.ID, result.Value, result.Tags)
}