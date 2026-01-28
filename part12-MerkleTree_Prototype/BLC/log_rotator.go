
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

type RotatingLogger struct {
	mu            sync.Mutex
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	logDir        string
	baseName      string
	rotationCount int
	maxRotations  int
}

func NewRotatingLogger(dir, name string, maxSize int64, maxRotations int) (*RotatingLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		maxSize:      maxSize,
		logDir:       dir,
		baseName:     name,
		maxRotations: maxRotations,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(rl.logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	oldPath := filepath.Join(rl.logDir, rl.baseName+".log")
	rotationTime := time.Now().Format("20060102_150405")
	archivePath := filepath.Join(rl.logDir, fmt.Sprintf("%s_%s.log.gz", rl.baseName, rotationTime))

	if err := compressFile(oldPath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(oldPath); err != nil {
		return err
	}

	rl.rotationCount++
	if rl.rotationCount > rl.maxRotations {
		rl.cleanupOldRotations()
	}

	return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (rl *RotatingLogger) cleanupOldRotations() {
	pattern := filepath.Join(rl.logDir, rl.baseName+"_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxRotations {
		filesToRemove := len(matches) - rl.maxRotations
		for i := 0; i < filesToRemove; i++ {
			os.Remove(matches[i])
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs", "app", 1024*1024, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}