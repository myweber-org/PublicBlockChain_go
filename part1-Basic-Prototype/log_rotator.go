package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    maxFiles    int
    currentSize int64
    file        *os.File
}

func NewRotator(filePath string, maxSize int64, maxFiles int) (*Rotator, error) {
    r := &Rotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
        return nil, err
    }

    if err := r.openCurrentFile(); err != nil {
        return nil, err
    }

    go r.timeBasedRotation()
    return r, nil
}

func (r *Rotator) openCurrentFile() error {
    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.file = file
    r.currentSize = stat.Size()
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) rotate() error {
    if r.file != nil {
        r.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }

    if err := r.cleanupOldFiles(); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }

    return r.openCurrentFile()
}

func (r *Rotator) cleanupOldFiles() error {
    pattern := fmt.Sprintf("%s.*", r.filePath)
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) <= r.maxFiles {
        return nil
    }

    for i := 0; i < len(files)-r.maxFiles; i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }
    return nil
}

func (r *Rotator) timeBasedRotation() {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        r.mu.Lock()
        if r.currentSize > 0 {
            if err := r.rotate(); err != nil {
                fmt.Printf("Time-based rotation failed: %v\n", err)
            }
        }
        r.mu.Unlock()
    }
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("/var/log/myapp/app.log", 10*1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}package main

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
	mu          sync.Mutex
	file        *os.File
	currentSize int64
	maxSize     int64
	basePath    string
	sequence    int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		maxSize:  maxSize,
		basePath: basePath,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	path := rl.basePath + ".log"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
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

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	oldPath := rl.basePath + ".log"
	archivePath := fmt.Sprintf("%s.%d.log.gz", rl.basePath, rl.sequence)
	rl.sequence++

	if err := compressFile(oldPath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(oldPath); err != nil {
		return err
	}

	return rl.openCurrent()
}

func compressFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app", 1024*1024)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
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

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    size       int64
    basePath   string
    currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
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

    n, err := rl.file.Write(p)
    if err != nil {
        return n, err
    }
    rl.size += int64(n)
    return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
    now := time.Now()
    dateStr := now.Format("2006-01-02")

    if rl.file == nil || rl.currentDay != dateStr || rl.size >= maxFileSize {
        if rl.file != nil {
            rl.file.Close()
            if err := rl.compressOldLog(); err != nil {
                return err
            }
            if err := rl.cleanupOldBackups(); err != nil {
                return err
            }
        }

        newPath := fmt.Sprintf("%s.%s.log", rl.basePath, dateStr)
        file, err := os.OpenFile(newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return err
        }

        info, err := file.Stat()
        if err != nil {
            file.Close()
            return err
        }

        rl.file = file
        rl.size = info.Size()
        rl.currentDay = dateStr
    }
    return nil
}

func (rl *RotatingLogger) compressOldLog() error {
    oldPath := fmt.Sprintf("%s.%s.log", rl.basePath, rl.currentDay)
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

func (rl *RotatingLogger) cleanupOldBackups() error {
    pattern := rl.basePath + ".*.log.gz"
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) > maxBackups {
        filesToDelete := files[:len(files)-maxBackups]
        for _, file := range filesToDelete {
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

    if rl.file != nil {
        return rl.file.Close()
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
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxFiles    = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	fileIndex   int
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: name,
	}

	if err := rl.openNextFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openNextFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	rl.fileIndex = (rl.fileIndex + 1) % maxFiles
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.fileIndex))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0

	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.openNextFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
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

	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	compression bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int, compression bool) (*RotatingLogger, error) {
	if maxSizeMB <= 0 {
		return nil, fmt.Errorf("maxSizeMB must be positive")
	}
	if maxFiles <= 0 {
		return nil, fmt.Errorf("maxFiles must be positive")
	}

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     int64(maxSizeMB) * 1024 * 1024,
		maxFiles:    maxFiles,
		compression: compression,
	}

	if err := rl.initialize(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) initialize() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	pattern := filepath.Base(rl.basePath) + ".*"
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return fmt.Errorf("failed to scan existing log files: %w", err)
	}

	rl.fileCount = len(matches)

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	rl.currentFile = file

	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	stat, err := rl.currentFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to stat log file: %w", err)
	}

	if stat.Size()+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log file: %w", err)
		}
	}

	return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

	if err := os.Rename(rl.basePath, archivePath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if rl.compression {
		go rl.compressFile(archivePath)
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		go rl.cleanupOldFiles()
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	rl.currentFile = file

	return nil
}

func (rl *RotatingLogger) compressFile(path string) {
	compressedPath := path + ".gz"
	if err := compressGzip(path, compressedPath); err != nil {
		log.Printf("Failed to compress %s: %v", path, err)
		return
	}
	if err := os.Remove(path); err != nil {
		log.Printf("Failed to remove uncompressed file %s: %v", path, err)
	}
}

func compressGzip(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gw := newGzipWriter(dstFile)
	defer gw.Close()

	_, err = io.Copy(gw, srcFile)
	return err
}

func (rl *RotatingLogger) cleanupOldFiles() {
	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)

	pattern := baseName + ".*"
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		log.Printf("Failed to scan log files for cleanup: %v", err)
		return
	}

	if len(matches) <= rl.maxFiles {
		return
	}

	filesToDelete := len(matches) - rl.maxFiles
	sortByModTime(matches)

	for i := 0; i < filesToDelete && i < len(matches); i++ {
		if err := os.Remove(matches[i]); err != nil {
			log.Printf("Failed to delete old log file %s: %v", matches[i], err)
		}
	}
}

func sortByModTime(files []string) {
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			infoI, _ := os.Stat(files[i])
			infoJ, _ := os.Stat(files[j])
			if infoI.ModTime().After(infoJ.ModTime()) {
				files[i], files[j] = files[j], files[i]
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 100; i++ {
		customLog.Printf("Log entry %d: Application is running normally", i)
		time.Sleep(100 * time.Millisecond)
	}
}