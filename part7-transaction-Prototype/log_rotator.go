
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
	currentNum int
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

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.size += int64(n)

	if err := rl.rotateIfNeeded(); err != nil {
		return n, err
	}
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.size < maxFileSize && rl.file != nil {
		return nil
	}

	if rl.file != nil {
		rl.file.Close()
		if err := rl.compressCurrent(); err != nil {
			return err
		}
		rl.cleanOldBackups()
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	rl.file = file
	info, err := file.Stat()
	if err != nil {
		return err
	}
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) compressCurrent() error {
	src, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer src.Close()

	backupName := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentNum)
	dst, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}
	rl.currentNum = (rl.currentNum + 1) % backupCount
	return os.Remove(rl.basePath)
}

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := fmt.Sprintf("%s.*.gz", rl.basePath)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	oldest := matches[0]
	for _, match := range matches[1:] {
		info1, _ := os.Stat(oldest)
		info2, _ := os.Stat(match)
		if info2.ModTime().Before(info1.ModTime()) {
			oldest = match
		}
	}
	os.Remove(oldest)
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
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
}