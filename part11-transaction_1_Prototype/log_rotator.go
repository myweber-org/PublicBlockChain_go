
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
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
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDay := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDay {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLog(); err != nil {
				log.Printf("Failed to compress old log: %v", err)
			}
			if err := rl.cleanupOldBackups(); err != nil {
				log.Printf("Failed to cleanup old backups: %v", err)
			}
		}

		newPath := rl.getLogPath(now)
		if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
			return fmt.Errorf("create directory failed: %w", err)
		}

		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("open log file failed: %w", err)
		}

		info, err := file.Stat()
		if err != nil {
			file.Close()
			return fmt.Errorf("stat log file failed: %w", err)
		}

		rl.file = file
		rl.size = info.Size()
		rl.currentDay = currentDay
	}
	return nil
}

func (rl *RotatingLogger) getLogPath(t time.Time) string {
	base := filepath.Base(rl.basePath)
	dir := filepath.Dir(rl.basePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", name, t.Format("2006-01-02"), ext))
}

func (rl *RotatingLogger) compressOldLog() error {
	oldPath := rl.getLogPath(time.Now().AddDate(0, 0, -1))
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil
	}

	compressedPath := oldPath + ".gz"
	inFile, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, inFile); err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	dir := filepath.Dir(rl.basePath)
	base := filepath.Base(rl.basePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	pattern := filepath.Join(dir, name+"-*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > backupCount {
		toDelete := matches[:len(matches)-backupCount]
		for _, file := range toDelete {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Application is running smoothly", i)
		time.Sleep(100 * time.Millisecond)
	}
}