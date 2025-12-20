package logrotator

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Rotator struct {
	mu            sync.Mutex
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	maxBackups    int
	compress      bool
	basePath      string
	rotationCount int
}

func NewRotator(basePath string, maxSize int64, maxBackups int, compress bool) (*Rotator, error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("maxSize must be positive")
	}

	r := &Rotator{
		maxSize:    maxSize,
		maxBackups: maxBackups,
		compress:   compress,
		basePath:   basePath,
	}

	if err := r.openCurrentFile(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentSize+int64(len(p)) > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := r.currentFile.Write(p)
	if err == nil {
		r.currentSize += int64(n)
	}
	return n, err
}

func (r *Rotator) rotate() error {
	if r.currentFile != nil {
		r.currentFile.Close()
	}

	r.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	oldPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

	if err := os.Rename(r.basePath, oldPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if r.compress {
		if err := r.compressFile(oldPath); err != nil {
			return fmt.Errorf("failed to compress log file: %w", err)
		}
		oldPath = oldPath + ".gz"
	}

	if err := r.cleanupOldBackups(); err != nil {
		return fmt.Errorf("failed to cleanup old backups: %w", err)
	}

	return r.openCurrentFile()
}

func (r *Rotator) compressFile(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(src + ".gz")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	return os.Remove(src)
}

func (r *Rotator) cleanupOldBackups() error {
	if r.maxBackups <= 0 {
		return nil
	}

	pattern := r.basePath + ".*"
	if r.compress {
		pattern += ".gz"
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= r.maxBackups {
		return nil
	}

	toDelete := matches[:len(matches)-r.maxBackups]
	for _, file := range toDelete {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (r *Rotator) openCurrentFile() error {
	file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	r.currentFile = file
	r.currentSize = info.Size()
	return nil
}

func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentFile != nil {
		return r.currentFile.Close()
	}
	return nil
}