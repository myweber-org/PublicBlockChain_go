
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    sequence    int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
        sequence: 0,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
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
    rotatedPath := fmt.Sprintf("%s.%s.%d", lr.basePath, timestamp, lr.sequence)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    lr.sequence++
    if lr.sequence > maxBackups {
        lr.cleanupOldBackups()
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

    os.Remove(source)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        filesToDelete := matches[:len(matches)-maxBackups]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()

    if lr.sequence == 0 {
        lr.determineSequence()
    }

    return nil
}

func (lr *LogRotator) determineSequence() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    maxSeq := 0
    for _, match := range matches {
        parts := strings.Split(filepath.Base(match), ".")
        if len(parts) >= 3 {
            seq, err := strconv.Atoi(parts[len(parts)-2])
            if err == nil && seq > maxSeq {
                maxSeq = seq
            }
        }
    }
    lr.sequence = maxSeq + 1
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }

        if i%100 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Println("Log rotation test completed")
}
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type RotatingWriter struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	backupCount int
}

func NewRotatingWriter(filePath string, maxSize int64, backupCount int) (*RotatingWriter, error) {
	writer := &RotatingWriter{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := writer.openCurrentFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *RotatingWriter) openCurrentFile() error {
	dir := filepath.Dir(w.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.currentFile = file
	w.currentSize = stat.Size()
	return nil
}

func (w *RotatingWriter) rotate() error {
	w.currentFile.Close()

	for i := w.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		if i == 0 {
			oldPath = w.filePath
		} else {
			oldPath = fmt.Sprintf("%s.%d", w.filePath, i)
		}
		newPath = fmt.Sprintf("%s.%d", w.filePath, i+1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	return w.openCurrentFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.currentFile.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.currentFile.Close()
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log", 1024*1024, 3)
	if err != nil {
		fmt.Printf("Failed to create rotating writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry number %d\n", i)
		if _, err := writer.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}