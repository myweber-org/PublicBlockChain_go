
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
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
}

func NewRotatingLogger(filePath string, maxSize int64, backupCount int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %w", err)
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %w", err)
	}

	rl.currentFile = file
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	info, err := rl.currentFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("stat file failed: %w", err)
	}

	if info.Size()+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("rotate failed: %w", err)
		}
	}

	return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return fmt.Errorf("close current file failed: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("rename file failed: %w", err)
	}

	if err := rl.compressBackup(backupPath); err != nil {
		return fmt.Errorf("compress backup failed: %w", err)
	}

	if err := rl.cleanOldBackups(); err != nil {
		return fmt.Errorf("clean old backups failed: %w", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressBackup(backupPath string) error {
	source, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("open backup file failed: %w", err)
	}
	defer source.Close()

	compressedPath := backupPath + ".gz"
	target, err := os.Create(compressedPath)
	if err != nil {
		return fmt.Errorf("create compressed file failed: %w", err)
	}
	defer target.Close()

	gzWriter := gzip.NewWriter(target)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, source); err != nil {
		return fmt.Errorf("compress data failed: %w", err)
	}

	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("remove original backup failed: %w", err)
	}

	return nil
}

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := rl.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob backup files failed: %w", err)
	}

	if len(matches) <= rl.backupCount {
		return nil
	}

	backupsToRemove := matches[:len(matches)-rl.backupCount]
	for _, backup := range backupsToRemove {
		if err := os.Remove(backup); err != nil {
			return fmt.Errorf("remove old backup %s failed: %w", backup, err)
		}
	}

	return nil
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
	logger, err := NewRotatingLogger("./logs/app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write failed: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize    = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logFileName   = "app.log"
)

type RotatingWriter struct {
	currentSize int64
	file        *os.File
}

func NewRotatingWriter() (*RotatingWriter, error) {
	w := &RotatingWriter{}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openFile() error {
	info, err := os.Stat(logFileName)
	if err == nil {
		w.currentSize = info.Size()
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	if w.currentSize+int64(len(p)) > maxLogSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.file.Write(p)
	w.currentSize += int64(n)
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	if err := os.Rename(logFileName, backupName); err != nil {
		return err
	}

	if err := w.openFile(); err != nil {
		return err
	}

	w.cleanupOldBackups()
	return nil
}

func (w *RotatingWriter) cleanupOldBackups() {
	pattern := fmt.Sprintf("%s.*", logFileName)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Failed to list backup files: %v", err)
		return
	}

	if len(matches) <= maxBackupFiles {
		return
	}

	oldestFiles := matches[:len(matches)-maxBackupFiles]
	for _, file := range oldestFiles {
		if err := os.Remove(file); err != nil {
			log.Printf("Failed to remove old backup %s: %v", file, err)
		}
	}
}

func (w *RotatingWriter) Close() error {
	return w.file.Close()
}

func main() {
	writer, err := NewRotatingWriter()
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, writer))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry number %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}