
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	mu          sync.Mutex
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentFile == nil || rl.currentSize >= maxFileSize {
		return rl.rotate()
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		if err := rl.currentFile.Close(); err != nil {
			return err
		}
		if err := rl.compressCurrentFile(); err != nil {
			log.Printf("Failed to compress log file: %v", err)
		}
		rl.cleanOldBackups()
	}

	timestamp := time.Now().Format("20060102_150405")
	newFileName := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

	file, err := os.Create(newFileName)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0

	return nil
}

func (rl *RotatingLogger) compressCurrentFile() error {
	if rl.currentFile == nil {
		return nil
	}

	oldPath := rl.currentFile.Name()
	compressedPath := oldPath + ".gz"

	src, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	if err := os.Remove(oldPath); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := filepath.Join(logDir, rl.baseName+"_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Failed to list backup files: %v", err)
		return
	}

	if len(matches) > backupCount {
		filesToRemove := matches[:len(matches)-backupCount]
		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				log.Printf("Failed to remove old backup %s: %v", file, err)
			}
		}
	}
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
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d: %s", i, "Sample log message for testing rotation")
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
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    currentFile *os.File
    maxSize     int64
    fileCount   int
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    dir := filepath.Dir(l.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    l.currentFile = file
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    info, err := l.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if info.Size()+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    return l.currentFile.Write(p)
}

func (l *RotatingLogger) rotate() error {
    if err := l.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", l.basePath, timestamp)
    if err := os.Rename(l.basePath, rotatedPath); err != nil {
        return err
    }

    compressedPath := rotatedPath + ".gz"
    if err := compressFile(rotatedPath, compressedPath); err != nil {
        return err
    }
    os.Remove(rotatedPath)

    l.fileCount++
    if l.fileCount > l.maxFiles {
        l.cleanupOldFiles()
    }

    return l.openCurrentFile()
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    gz := gzip.NewWriter(destination)
    defer gz.Close()

    _, err = io.Copy(gz, source)
    return err
}

func (l *RotatingLogger) cleanupOldFiles() {
    dir := filepath.Dir(l.basePath)
    baseName := filepath.Base(l.basePath)

    files, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var compressedFiles []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            compressedFiles = append(compressedFiles, filepath.Join(dir, name))
        }
    }

    if len(compressedFiles) > l.maxFiles {
        filesToDelete := compressedFiles[:len(compressedFiles)-l.maxFiles]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (l *RotatingLogger) extractTimestamp(filename string) (time.Time, error) {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return time.Time{}, fmt.Errorf("invalid filename format")
    }
    timestampStr := parts[len(parts)-2]
    return time.Parse("20060102_150405", timestampStr)
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed. Check /var/log/myapp directory for rotated files.")
}