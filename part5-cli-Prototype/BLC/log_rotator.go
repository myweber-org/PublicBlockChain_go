package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize   = 1024 * 1024 // 1MB
	maxArchives  = 5
	logFileName  = "app.log"
	archiveDir   = "archives"
)

func rotateLog() error {
	info, err := os.Stat(logFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() < maxLogSize {
		return nil
	}

	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveName := filepath.Join(archiveDir, fmt.Sprintf("%s_%s.log", logFileName, timestamp))

	if err := os.Rename(logFileName, archiveName); err != nil {
		return err
	}

	cleanupOldArchives()
	return nil
}

func cleanupOldArchives() {
	files, err := filepath.Glob(filepath.Join(archiveDir, "*.log"))
	if err != nil {
		return
	}

	if len(files) <= maxArchives {
		return
	}

	for i := 0; i < len(files)-maxArchives; i++ {
		os.Remove(files[i])
	}
}

func writeLog(message string) error {
	if err := rotateLog(); err != nil {
		return err
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	_, err = io.WriteString(file, logEntry)
	return err
}

func main() {
	for i := 1; i <= 100; i++ {
		message := fmt.Sprintf("Log entry number %d", i)
		if err := writeLog(message); err != nil {
			fmt.Printf("Error writing log: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("Log rotation demonstration completed")
}