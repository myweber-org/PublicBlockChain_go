
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type RotatingLogger struct {
	mu          sync.Mutex
	file        *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	backupCount int
}

func NewRotatingLogger(filePath string, maxSize int64, backupCount int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()

	for i := rl.backupCount - 1; i >= 0; i-- {
		src := rl.backupPath(i)
		dst := rl.backupPath(i + 1)

		if _, err := os.Stat(src); err == nil {
			if err := os.Rename(src, dst); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(rl.filePath, rl.backupPath(0)); err != nil {
		return err
	}

	return rl.openFile()
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.filePath + ".1"
	}
	return fmt.Sprintf("%s.%d", rl.filePath, index+1)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 1024*1024, 3)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d: This is a test log message\n", i)
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
	}

	fmt.Println("Log rotation test completed")
}