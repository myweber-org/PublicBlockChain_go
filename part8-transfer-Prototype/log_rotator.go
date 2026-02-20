
package main

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
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    currentFile   *os.File
    currentSize   int64
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }

    err := logger.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return err
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
        err := rl.rotate()
        if err != nil {
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
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    err := os.Rename(rl.basePath, rotatedPath)
    if err != nil {
        return err
    }

    rl.rotationCount++
    go rl.compressOldLog(rotatedPath)

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressOldLog(path string) {
    compressedPath := path + ".gz"

    srcFile, err := os.Open(path)
    if err != nil {
        return
    }
    defer srcFile.Close()

    destFile, err := os.Create(compressedPath)
    if err != nil {
        return
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return
    }

    os.Remove(path)
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func (rl *RotatingLogger) ListArchives() []string {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return []string{}
    }

    var archives []string
    for _, match := range matches {
        archives = append(archives, filepath.Base(match))
    }
    return archives
}

func (rl *RotatingLogger) CleanOldArchives(maxAgeDays int) {
    cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
    pattern := rl.basePath + ".*.gz"

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    for _, match := range matches {
        info, err := os.Stat(match)
        if err != nil {
            continue
        }

        if info.ModTime().Before(cutoff) {
            os.Remove(match)
        }
    }
}

func main() {
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        os.Exit(1)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Iteration %d: Processing data chunk %d\n",
            time.Now().Format(time.RFC3339),
            i,
            i*1024)

        _, err := logger.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Write error: %v\n", err)
        }

        time.Sleep(10 * time.Millisecond)
    }

    archives := logger.ListArchives()
    fmt.Printf("Created %d archive(s): %v\n", len(archives), archives)

    logger.CleanOldArchives(7)
    fmt.Println("Cleaned archives older than 7 days")
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

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(filePath string, maxSizeMB int) (*RotatingLogger, error) {
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

	return &RotatingLogger{
		currentFile:  file,
		filePath:     filePath,
		maxSize:      maxSize,
		currentSize:  info.Size(),
		rotationCount: 0,
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

	rl.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.filePath, timestamp)

	if err := rl.compressFile(rl.filePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(rl.filePath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

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
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application event occurred at %v\n", i, time.Now())
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}