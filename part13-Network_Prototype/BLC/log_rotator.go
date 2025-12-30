
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	backupCount int
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) *LogRotator {
	return &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if err := lr.rotateIfNeeded(); err != nil {
		return 0, err
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write(p)
}

func (lr *LogRotator) rotateIfNeeded() error {
	info, err := os.Stat(lr.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if info.Size() < lr.maxSize {
		return nil
	}

	for i := lr.backupCount - 1; i > 0; i-- {
		oldName := fmt.Sprintf("%s.%d", lr.filePath, i)
		newName := fmt.Sprintf("%s.%d", lr.filePath, i+1)
		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	backupName := fmt.Sprintf("%s.1", lr.filePath)
	if err := os.Rename(lr.filePath, backupName); err != nil {
		return err
	}

	return nil
}

func (lr *LogRotator) CleanOldBackups() error {
	for i := lr.backupCount + 1; ; i++ {
		backupPath := fmt.Sprintf("%s.%d", lr.filePath, i)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		if err := os.Remove(backupPath); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	rotator := NewLogRotator("app.log", 1024*1024, 5)

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		rotator.Write([]byte(logEntry))
	}

	rotator.CleanOldBackups()
	fmt.Println("Log rotation completed")
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"
)

const (
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logDir         = "./logs"
)

type LogRotator struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewLogRotator(baseName string) (*LogRotator, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    filePath := filepath.Join(logDir, baseName+".log")
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        currentFile:  file,
        filePath:     filePath,
        bytesWritten: info.Size(),
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.bytesWritten+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.bytesWritten += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := strings.TrimSuffix(lr.filePath, ".log") + "_" + timestamp + ".log"

    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = file
    lr.bytesWritten = 0

    go lr.cleanupOldFiles()
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := filepath.Join(logDir, strings.TrimSuffix(filepath.Base(lr.filePath), ".log") + "_*.log")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupFiles {
        return
    }

    sort.Strings(matches)
    filesToDelete := matches[:len(matches)-maxBackupFiles]

    for _, file := range filesToDelete {
        os.Remove(file)
    }
}

func (lr *LogRotator) Close() error {
    return lr.currentFile.Close()
}

func main() {
    rotator, err := NewLogRotator("application")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}