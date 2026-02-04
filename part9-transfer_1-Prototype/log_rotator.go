
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
	fileCounter  int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     int64(maxSizeMB) * 1024 * 1024,
		fileCounter: 1,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileCounter)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
		go rl.compressOldFile(rl.fileCounter)
		rl.fileCounter++
		return rl.openCurrentFile()
	}
	return nil
}

func (rl *RotatingLogger) compressOldFile(fileNum int) {
	oldName := fmt.Sprintf("%s.%d.log", rl.basePath, fileNum)
	compressedName := fmt.Sprintf("%s.%d.log.gz", rl.basePath, fileNum)

	oldFile, err := os.Open(oldName)
	if err != nil {
		return
	}
	defer oldFile.Close()

	compressedFile, err := os.Create(compressedName)
	if err != nil {
		return
	}
	defer compressedFile.Close()

	gzWriter := gzip.NewWriter(compressedFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		os.Remove(compressedName)
		return
	}

	os.Remove(oldName)
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
	logger, err := NewRotatingLogger("application", 10)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Processing request from client\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogFiles = 5
    logFileName = "app.log"
)

func rotateLogs() error {
    for i := maxLogFiles - 1; i > 0; i-- {
        oldName := fmt.Sprintf("%s.%d", logFileName, i)
        newName := fmt.Sprintf("%s.%d", logFileName, i+1)

        if _, err := os.Stat(oldName); err == nil {
            err := os.Rename(oldName, newName)
            if err != nil {
                return fmt.Errorf("failed to rename %s to %s: %w", oldName, newName, err)
            }
        }
    }

    if _, err := os.Stat(logFileName); err == nil {
        err := os.Rename(logFileName, fmt.Sprintf("%s.1", logFileName))
        if err != nil {
            return fmt.Errorf("failed to rename current log: %w", err)
        }
    }

    return nil
}

func writeLog(message string) error {
    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

    _, err = file.WriteString(logEntry)
    return err
}

func main() {
    fileInfo, err := os.Stat(logFileName)
    if err == nil && fileInfo.Size() > 1024*1024 {
        fmt.Println("Log file exceeds 1MB, rotating...")
        if err := rotateLogs(); err != nil {
            fmt.Printf("Log rotation failed: %v\n", err)
            return
        }
    }

    if err := writeLog("Application started"); err != nil {
        fmt.Printf("Failed to write log: %v\n", err)
        return
    }

    fmt.Println("Log operation completed successfully")
}