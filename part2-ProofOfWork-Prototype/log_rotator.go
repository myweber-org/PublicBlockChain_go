
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxBackups   = 5
	logExtension = ".log"
)

type LogRotator struct {
	currentFile *os.File
	filePath    string
	baseName    string
	dir         string
	written     int64
}

func NewLogRotator(filePath string) (*LogRotator, error) {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	baseName := strings.TrimSuffix(base, filepath.Ext(base))

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		currentFile: file,
		filePath:    filePath,
		baseName:    baseName,
		dir:         dir,
		written:     info.Size(),
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.written+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.written += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(lr.dir, fmt.Sprintf("%s_%s%s", lr.baseName, timestamp, logExtension))
	if err := os.Rename(lr.filePath, oldPath); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.written = 0

	go lr.cleanupOldFiles()
	return nil
}

func (lr *LogRotator) cleanupOldFiles() {
	pattern := filepath.Join(lr.dir, lr.baseName+"_*"+logExtension)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackups {
		return
	}

	filesToRemove := matches[:len(matches)-maxBackups]
	for _, file := range filesToRemove {
		os.Remove(file)
	}
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("./logs/application.log")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}