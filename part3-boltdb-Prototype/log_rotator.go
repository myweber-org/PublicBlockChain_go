package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRotator struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
	maxRotations int
}

func NewLogRotator(filePath string, maxSize int64, maxRotations int) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath:     filePath,
		maxSize:      maxSize,
		maxRotations: maxRotations,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedFile := fmt.Sprintf("%s.%s", lr.filePath, timestamp)
	
	if err := os.Rename(lr.filePath, rotatedFile); err != nil {
		return err
	}

	if err := lr.compressFile(rotatedFile); err != nil {
		return err
	}

	lr.rotationCount++
	if lr.rotationCount > lr.maxRotations {
		lr.cleanupOldRotations()
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(source string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(source + ".gz")
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	if err != nil {
		return err
	}

	return os.Remove(source)
}

func (lr *LogRotator) openCurrentFile() error {
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) cleanupOldRotations() {
	files, err := filepath.Glob(lr.filePath + ".*.gz")
	if err != nil {
		return
	}

	if len(files) > lr.maxRotations {
		filesToDelete := files[:len(files)-lr.maxRotations]
		for _, file := range filesToDelete {
			os.Remove(file)
		}
	}
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}