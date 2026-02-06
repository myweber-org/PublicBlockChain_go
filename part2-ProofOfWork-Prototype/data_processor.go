
package data_processor

import (
	"regexp"
	"strings"
)

type Processor struct {
	allowedPattern *regexp.Regexp
}

func NewProcessor(allowedPattern string) (*Processor, error) {
	compiled, err := regexp.Compile(allowedPattern)
	if err != nil {
		return nil, err
	}
	return &Processor{allowedPattern: compiled}, nil
}

func (p *Processor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return p.allowedPattern.FindString(trimmed)
}

func (p *Processor) Validate(input string) bool {
	return p.allowedPattern.MatchString(input)
}

func (p *Processor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := p.CleanInput(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}