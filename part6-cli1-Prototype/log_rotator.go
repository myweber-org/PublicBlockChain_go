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
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	filePath    string
	mu          sync.Mutex
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: basePath,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
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

	timestamp := time.Now().Format("20060102150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.filePath, timestamp)

	if err := compressFile(rl.filePath, archivePath); err != nil {
		return err
	}

	if err := cleanupOldFiles(rl.filePath); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func compressFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func cleanupOldFiles(basePath string) error {
	pattern := basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		filesToRemove := matches[:len(matches)-maxBackups]
		for _, file := range filesToRemove {
			os.Remove(file)
		}
	}
	return nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = stat.Size()
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
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
    maxBackupLogs = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentSize int64
    file        *os.File
}

func NewLogRotator() (*LogRotator, error) {
    lr := &LogRotator{}
    if err := lr.openLogFile(); err != nil {
        return nil, err
    }
    return lr, nil
}

func (lr *LogRotator) openLogFile() error {
    info, err := os.Stat(logFileName)
    if err != nil && !os.IsNotExist(err) {
        return err
    }
    if err == nil {
        lr.currentSize = info.Size()
    }

    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    lr.file = file
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxLogSize {
        if err := lr.rotate(); err != nil {
            log.Printf("Failed to rotate log: %v", err)
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    if err := os.Rename(logFileName, backupName); err != nil {
        return err
    }

    if err := lr.openLogFile(); err != nil {
        return err
    }
    lr.currentSize = 0

    go lr.cleanupOldLogs()
    return nil
}

func (lr *LogRotator) cleanupOldLogs() {
    files, err := filepath.Glob(logFileName + ".*")
    if err != nil {
        return
    }

    if len(files) <= maxBackupLogs {
        return
    }

    for i := 0; i < len(files)-maxBackupLogs; i++ {
        os.Remove(files[i])
    }
}

func (lr *LogRotator) Close() error {
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator()
    if err != nil {
        log.Fatal(err)
    }
    defer rotator.Close()

    logger := log.New(io.MultiWriter(os.Stdout, rotator), "", log.LstdFlags)

    for i := 0; i < 100; i++ {
        logger.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(100 * time.Millisecond)
    }
}