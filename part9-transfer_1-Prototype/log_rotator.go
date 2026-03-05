
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
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
        basePath:    basePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
        sequence:    0,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) rotate() error {
    if rl.file != nil {
        rl.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.%d", rl.basePath, timestamp, rl.sequence)

    if err := os.Rename(rl.basePath, archivedPath); err != nil {
        return err
    }

    if err := rl.compressFile(archivedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    rl.sequence++

    return nil
}

func (rl *RotatingLogger) compressFile(source string) error {
    reader, err := os.Open(source)
    if err != nil {
        return err
    }
    defer reader.Close()

    compressedPath := source + ".gz"
    writer, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer writer.Close()

    gzWriter := gzip.NewWriter(writer)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, reader); err != nil {
        return err
    }

    os.Remove(source)
    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles(maxFiles int) error {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxFiles {
        return nil
    }

    for i := 0; i < len(matches)-maxFiles; i++ {
        if err := os.Remove(matches[i]); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    if err := logger.cleanupOldFiles(5); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }

    fmt.Println("Log rotation test completed")
}package main

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

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
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
	oldPath := filepath.Join(logDir, rl.baseName+".log")
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", rl.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.compressFile(newPath); err != nil {
		log.Printf("Compression failed for %s: %v", newPath, err)
	}

	if err := rl.cleanupOldFiles(); err != nil {
		log.Printf("Cleanup failed: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dst := src + ".gz"
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gz := gzip.NewWriter(dstFile)
	defer gz.Close()

	if _, err := io.Copy(gz, srcFile); err != nil {
		return err
	}

	return os.Remove(src)
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := filepath.Join(logDir, rl.baseName+".*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > backupCount {
		filesToRemove := matches[:len(matches)-backupCount]
		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("application")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			log.Printf("Write failed: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}