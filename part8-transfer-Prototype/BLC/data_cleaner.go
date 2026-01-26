package csvutils

import (
	"encoding/csv"
	"io"
	"strings"
	"unicode"
)

type Cleaner struct {
	TrimSpaces  bool
	RemoveEmpty bool
	ToLowercase bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		TrimSpaces:  true,
		RemoveEmpty: true,
		ToLowercase: false,
	}
}

func (c *Cleaner) CleanRecord(record []string) []string {
	var cleaned []string

	for _, field := range record {
		processed := field

		if c.TrimSpaces {
			processed = strings.TrimSpace(processed)
		}

		if c.ToLowercase {
			processed = strings.ToLower(processed)
		}

		processed = strings.Map(func(r rune) rune {
			if unicode.IsControl(r) && r != '\n' && r != '\t' {
				return -1
			}
			return r
		}, processed)

		if !c.RemoveEmpty || processed != "" {
			cleaned = append(cleaned, processed)
		}
	}

	return cleaned
}

func (c *Cleaner) ProcessCSV(reader *csv.Reader, writer *csv.Writer) error {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleaned := c.CleanRecord(record)
		if len(cleaned) > 0 {
			if err := writer.Write(cleaned); err != nil {
				return err
			}
		}
	}

	return writer.Flush()
}