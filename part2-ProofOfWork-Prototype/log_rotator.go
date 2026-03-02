
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
		sequence: 0,
	}

	if err := rl.openNewFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
		rl.compressOldFile()
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%03d.log", rl.baseName, timestamp, rl.sequence)
	path := filepath.Join(logDir, filename)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	rl.sequence++
	return nil
}

func (rl *RotatingLogger) compressOldFile() {
	if rl.sequence <= 1 {
		return
	}

	oldSequence := rl.sequence - 2
	timestamp := time.Now().Add(-time.Minute).Format("20060102_150405")
	oldFilename := fmt.Sprintf("%s_%s_%03d.log", rl.baseName, timestamp, oldSequence)
	oldPath := filepath.Join(logDir, oldFilename)
	compressedPath := oldPath + ".gz"

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return
	}

	source, err := os.Open(oldPath)
	if err != nil {
		return
	}
	defer source.Close()

	dest, err := os.Create(compressedPath)
	if err != nil {
		return
	}
	defer dest.Close()

	gzWriter := gzip.NewWriter(dest)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, source); err != nil {
		return
	}

	os.Remove(oldPath)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.openNewFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}