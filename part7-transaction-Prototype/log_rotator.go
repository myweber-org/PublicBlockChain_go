
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
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
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
		rl.currentFile = nil
	}

	rl.rotationCount++
	archiveName := fmt.Sprintf("%s.%d.%s.gz", 
		rl.basePath, 
		rl.rotationCount, 
		time.Now().Format("20060102_150405"))

	sourcePath := rl.basePath
	if err := compressFile(sourcePath, archiveName); err != nil {
		return fmt.Errorf("compression failed: %v", err)
	}

	if err := os.Remove(sourcePath); err != nil {
		return fmt.Errorf("failed to remove old log: %v", err)
	}

	return rl.openCurrentFile()
}

func compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
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
	logger, err := NewRotatingLogger("/var/log/app/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write failed: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024
	maxBackups   = 5
	logDirectory = "./logs"
)

type RotatingLogger struct {
	filename   string
	current    *os.File
	size       int64
	mu         sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDirectory, 0755); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(logDirectory, filename)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		filename: filename,
		current:  file,
		size:     info.Size(),
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.current.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.current != nil {
		rl.current.Close()
	}

	baseName := filepath.Join(logDirectory, rl.filename)
	for i := maxBackups - 1; i >= 0; i-- {
		oldName := fmt.Sprintf("%s.%d", baseName, i)
		newName := fmt.Sprintf("%s.%d", baseName, i+1)

		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	if _, err := os.Stat(baseName); err == nil {
		os.Rename(baseName, baseName+".0")
	}

	file, err := os.OpenFile(baseName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.current = file
	rl.size = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
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

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        err := lr.rotate()
        if err != nil {
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
        err := lr.compressCurrentFile()
        if err != nil {
            return err
        }
    }

    lr.sequence++
    if lr.sequence > maxBackups {
        lr.cleanOldBackups()
    }

    return lr.openCurrentFile()
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
    return nil
}

func (lr *LogRotator) compressCurrentFile() error {
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s.gz", lr.basePath, timestamp)

    source, err := os.Open(lr.basePath)
    if err != nil {
        return err
    }
    defer source.Close()

    dest, err := os.Create(backupPath)
    if err != nil {
        return err
    }
    defer dest.Close()

    gzWriter := gzip.NewWriter(dest)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, source)
    if err != nil {
        return err
    }

    return os.Remove(lr.basePath)
}

func (lr *LogRotator) cleanOldBackups() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    backupMap := make(map[string]time.Time)
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        timestampStr := parts[len(parts)-2]
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            continue
        }
        backupMap[match] = t
    }

    var sortedBackups []string
    for backup := range backupMap {
        sortedBackups = append(sortedBackups, backup)
    }

    for i := 0; i < len(sortedBackups)-maxBackups; i++ {
        os.Remove(sortedBackups[i])
    }
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
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
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
    filename   string
    current    *os.File
    currentSize int64
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    if err := rotator.openCurrent(); err != nil {
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

    n, err := lr.current.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.current != nil {
        lr.current.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedFile := fmt.Sprintf("%s.%s", lr.filename, timestamp)
    if err := os.Rename(lr.filename, rotatedFile); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedFile); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openCurrent()
}

func (lr *LogRotator) compressFile(filename string) error {
    src, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(filename + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(filename)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    var timestamps []int64
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        ts, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
        if err != nil {
            continue
        }
        timestamps = append(timestamps, ts)
    }

    for i := 0; i < len(timestamps)-maxBackups; i++ {
        oldestIdx := i
        for j := i + 1; j < len(timestamps); j++ {
            if timestamps[j] < timestamps[oldestIdx] {
                oldestIdx = j
            }
        }
        os.Remove(matches[oldestIdx])
        timestamps[oldestIdx] = timestamps[i]
        matches[oldestIdx] = matches[i]
    }

    return nil
}

func (lr *LogRotator) openCurrent() error {
    file, err := os.OpenFile(lr.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.current = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.current != nil {
        return lr.current.Close()
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
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
}