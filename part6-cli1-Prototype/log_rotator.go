
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
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	}

	rl.rotationCount++
	archiveName := fmt.Sprintf("%s.%d.%s.gz", 
		rl.basePath, 
		rl.rotationCount, 
		time.Now().Format("20060102_150405"))

	if err := rl.compressFile(rl.basePath, archiveName); err != nil {
		return err
	}

	if err := os.Truncate(rl.basePath, 0); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
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
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

const (
    maxSize    = 10 * 1024 * 1024 // 10MB
    maxBackups = 5
)

type RotatingLog struct {
    file     *os.File
    size     int64
    basePath string
}

func NewRotatingLog(path string) (*RotatingLog, error) {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLog{
        file:     file,
        size:     info.Size(),
        basePath: path,
    }, nil
}

func (r *RotatingLog) Write(p []byte) (int, error) {
    if r.size+int64(len(p)) > maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    r.size += int64(n)
    return n, err
}

func (r *RotatingLog) rotate() error {
    if err := r.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    if err := os.Rename(r.basePath, rotatedPath); err != nil {
        return err
    }

    if err := r.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(r.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    r.file = file
    r.size = 0
    r.cleanupOldBackups()
    return nil
}

func (r *RotatingLog) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(path)
    return nil
}

func (r *RotatingLog) cleanupOldBackups() {
    pattern := r.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
}

func (r *RotatingLog) Close() error {
    return r.file.Close()
}

func main() {
    log, err := NewRotatingLog("application.log")
    if err != nil {
        panic(err)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingFile struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
}

func NewRotatingFile(filePath string, maxSize int64) (*RotatingFile, error) {
    rf := &RotatingFile{
        filePath: filePath,
        maxSize:  maxSize,
    }

    if err := rf.openOrCreate(); err != nil {
        return nil, err
    }

    return rf, nil
}

func (rf *RotatingFile) openOrCreate() error {
    info, err := os.Stat(rf.filePath)
    if os.IsNotExist(err) {
        dir := filepath.Dir(rf.filePath)
        if err := os.MkdirAll(dir, 0755); err != nil {
            return err
        }
        file, err := os.Create(rf.filePath)
        if err != nil {
            return err
        }
        rf.file = file
        rf.currentSize = 0
        return nil
    } else if err != nil {
        return err
    }

    file, err := os.OpenFile(rf.filePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    rf.file = file
    rf.currentSize = info.Size()
    return nil
}

func (rf *RotatingFile) rotate() error {
    if err := rf.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    baseName := strings.TrimSuffix(rf.filePath, filepath.Ext(rf.filePath))
    ext := filepath.Ext(rf.filePath)
    rotatedPath := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)

    if err := os.Rename(rf.filePath, rotatedPath); err != nil {
        return err
    }

    return rf.openOrCreate()
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
    rf.mu.Lock()
    defer rf.mu.Unlock()

    if rf.currentSize+int64(len(p)) > rf.maxSize {
        if err := rf.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rf.file.Write(p)
    if err == nil {
        rf.currentSize += int64(n)
    }
    return n, err
}

func (rf *RotatingFile) Close() error {
    rf.mu.Lock()
    defer rf.mu.Unlock()
    if rf.file != nil {
        return rf.file.Close()
    }
    return nil
}

func main() {
    rf, err := NewRotatingFile("./logs/app.log", 1024*1024)
    if err != nil {
        fmt.Printf("Failed to create rotating file: %v\n", err)
        return
    }
    defer rf.Close()

    for i := 0; i < 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rf.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}