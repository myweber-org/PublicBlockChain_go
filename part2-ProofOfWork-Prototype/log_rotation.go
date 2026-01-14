package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	filePath    string
	maxSize     int64
	maxAge      time.Duration
	currentFile *os.File
	logger      *log.Logger
	written     int64
}

func NewRotatingLogger(filePath string, maxSize int64, maxAge time.Duration) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: filePath,
		maxSize:  maxSize,
		maxAge:   maxAge,
	}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	go rl.cleanupOldFiles()
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.written = info.Size()
	rl.logger = log.New(file, "", log.LstdFlags)
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.written+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.written += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := rl.filePath + "." + timestamp

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	return rl.openFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		files, err := filepath.Glob(rl.filePath + ".*")
		if err != nil {
			continue
		}

		cutoff := time.Now().Add(-rl.maxAge)
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				os.Remove(file)
			}
		}
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu         sync.Mutex
    basePath   string
    maxSize    int64
    maxFiles   int
    current    *os.File
    currentSize int64
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openCurrent(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (l *RotatingLogger) openCurrent() error {
    file, err := os.OpenFile(l.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.current = file
    l.currentSize = info.Size()
    return nil
}

func (l *RotatingLogger) rotate() error {
    l.current.Close()

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", l.basePath, timestamp)

    if err := os.Rename(l.basePath, rotatedPath); err != nil {
        return err
    }

    if err := l.openCurrent(); err != nil {
        return err
    }

    l.cleanupOldFiles()
    return nil
}

func (l *RotatingLogger) cleanupOldFiles() {
    dir := filepath.Dir(l.basePath)
    base := filepath.Base(l.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var logFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, base+".") && !entry.IsDir() {
            logFiles = append(logFiles, filepath.Join(dir, name))
        }
    }

    if len(logFiles) <= l.maxFiles {
        return
    }

    sort.Strings(logFiles)
    filesToRemove := logFiles[:len(logFiles)-l.maxFiles]

    for _, file := range filesToRemove {
        os.Remove(file)
    }
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := l.current.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    if l.current != nil {
        return l.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}