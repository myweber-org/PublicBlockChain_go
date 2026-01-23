
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingWriter struct {
	mu         sync.Mutex
	file       *os.File
	filePath   string
	currentSize int64
}

func NewRotatingWriter(filePath string) (*RotatingWriter, error) {
	writer := &RotatingWriter{
		filePath: filePath,
	}

	if err := writer.openFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *RotatingWriter) openFile() error {
	dir := filepath.Dir(w.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.file = file
	w.currentSize = stat.Size()
	return nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > maxFileSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if w.file != nil {
		w.file.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", w.filePath, timestamp)

	if err := os.Rename(w.filePath, backupPath); err != nil {
		return err
	}

	if err := w.cleanupOldBackups(); err != nil {
		fmt.Printf("Warning: Failed to cleanup old backups: %v\n", err)
	}

	return w.openFile()
}

func (w *RotatingWriter) cleanupOldBackups() error {
	pattern := fmt.Sprintf("%s.*", w.filePath)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackups {
		return nil
	}

	oldestFiles := matches[:len(matches)-maxBackups]
	for _, file := range oldestFiles {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log")
	if err != nil {
		fmt.Printf("Failed to create rotating writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := writer.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type RotatingLogger struct {
    currentFile   *os.File
    basePath      string
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentFile:   file,
        basePath:      basePath,
        maxSize:       maxSize,
        currentSize:   info.Size(),
        rotationCount: 0,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    if err := compressFile(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    rl.rotationCount++
    return nil
}

func compressFile(source, target string) error {
    sourceFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    targetFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer targetFile.Close()

    gzWriter := gzip.NewWriter(targetFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    return err
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Printf("Total rotations: %d\n", logger.rotationCount)
}