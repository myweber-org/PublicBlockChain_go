
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
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rl.currentFile = f
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.rotationCount, timestamp)
	rl.rotationCount++

	if err := compressFile(rl.basePath, archivePath); err != nil {
		return fmt.Errorf("compression failed: %w", err)
	}

	if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	_, err = io.Copy(gz, source)
	return err
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
	logger, err := NewRotatingLogger("/var/log/app/app.log", 10)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    fileIndex   int
}

func NewRotatingWriter(basePath string) (*RotatingWriter, error) {
    w := &RotatingWriter{
        basePath: basePath,
    }
    if err := w.openCurrentFile(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
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

func (w *RotatingWriter) rotate() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

    w.fileIndex++
    if w.fileIndex > maxBackups {
        w.fileIndex = 1
    }

    oldPath := fmt.Sprintf("%s.%d", w.basePath, w.fileIndex)
    if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
        return err
    }

    if err := os.Rename(w.basePath, oldPath); err != nil && !os.IsNotExist(err) {
        return err
    }

    return w.openCurrentFile()
}

func (w *RotatingWriter) openCurrentFile() error {
    file, err := os.OpenFile(w.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    w.currentFile = file
    w.currentSize = info.Size()
    return nil
}

func (w *RotatingWriter) Close() error {
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        return
    }
    defer writer.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}