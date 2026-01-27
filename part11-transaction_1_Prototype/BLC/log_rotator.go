
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
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	rl := &RotatingLogger{
		baseName: name,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	rl.currentFile = file
	info, _ := file.Stat()
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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

	oldLogs, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return fmt.Errorf("failed to list log files: %w", err)
	}

	if len(oldLogs) >= maxBackups {
		if err := rl.compressAndRemove(oldLogs[0]); err != nil {
			return err
		}
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressAndRemove(filename string) error {
	src, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open log file for compression: %w", err)
	}
	defer src.Close()

	dest, err := os.Create(filename + ".gz")
	if err != nil {
		return fmt.Errorf("failed to create compressed file: %w", err)
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return fmt.Errorf("failed to compress log file: %w", err)
	}

	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to remove original log file: %w", err)
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
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

type RotatingLog struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    rotationNum int
}

func NewRotatingLog(filePath string, maxSizeMB int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLog{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
        rotationNum: 0,
    }, nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
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

func (rl *RotatingLog) rotate() error {
    if rl.file != nil {
        rl.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedName := fmt.Sprintf("%s.%s.%d.gz", rl.filePath, timestamp, rl.rotationNum)

    if err := rl.compressCurrentLog(archivedName); err != nil {
        return err
    }

    if err := os.Truncate(rl.filePath, 0); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    rl.rotationNum++

    return nil
}

func (rl *RotatingLog) compressCurrentLog(destPath string) error {
    src, err := os.Open(rl.filePath)
    if err != nil {
        return err
    }
    defer src.Close()

    dest, err := os.Create(destPath)
    if err != nil {
        return err
    }
    defer dest.Close()

    gz := gzip.NewWriter(dest)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    return err
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func (rl *RotatingLog) CleanupOldLogs(maxAgeDays int) error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    cutoffTime := time.Now().AddDate(0, 0, -maxAgeDays)
    pattern := filepath.Join(filepath.Dir(rl.filePath), filepath.Base(rl.filePath)+".*.gz")

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    for _, match := range matches {
        info, err := os.Stat(match)
        if err != nil {
            continue
        }

        if info.ModTime().Before(cutoffTime) {
            os.Remove(match)
        }
    }

    return nil
}