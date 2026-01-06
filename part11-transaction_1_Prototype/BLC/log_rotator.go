
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.size += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	if rl.file == nil || rl.currentDay != currentDate || rl.size >= maxFileSize {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLog(); err != nil {
				return err
			}
			if err := rl.cleanupOldBackups(); err != nil {
				return err
			}
		}

		rl.currentDay = currentDate
		filename := fmt.Sprintf("%s.%s.log", rl.basePath, currentDate)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		info, err := file.Stat()
		if err != nil {
			file.Close()
			return err
		}

		rl.file = file
		rl.size = info.Size()
	}
	return nil
}

func (rl *RotatingLogger) compressOldLog() error {
	oldPath := fmt.Sprintf("%s.%s.log", rl.basePath, rl.currentDay)
	compressedPath := oldPath + ".gz"

	src, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	files, err := filepath.Glob(rl.basePath + ".*.log.gz")
	if err != nil {
		return err
	}

	if len(files) > backupCount {
		sortFilesByModTime(files)
		for i := 0; i < len(files)-backupCount; i++ {
			os.Remove(files[i])
		}
	}
	return nil
}

func sortFilesByModTime(files []string) {
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			infoI, _ := os.Stat(files[i])
			infoJ, _ := os.Stat(files[j])
			if infoI.ModTime().After(infoJ.ModTime()) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}