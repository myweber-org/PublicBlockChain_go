package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxBackupCount = 5
	logFileName   = "app.log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fullPath := filepath.Join(basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		basePath:    basePath,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("rotation failed: %w", err)
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	oldPath := filepath.Join(lr.basePath, logFileName)
	newPath := filepath.Join(lr.basePath, backupName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	lr.currentFile = file
	lr.currentSize = 0

	go lr.cleanupOldBackups()

	return nil
}

func (lr *LogRotator) cleanupOldBackups() {
	pattern := filepath.Join(lr.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackupCount {
		return
	}

	sort.Strings(matches)
	filesToRemove := matches[:len(matches)-maxBackupCount]

	for _, file := range filesToRemove {
		os.Remove(file)
	}
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("./logs")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		}
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

    if err := rotator.openCurrentFile(); err != nil {
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

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
        if err := lr.compressCurrentFile(); err != nil {
            return err
        }
        lr.cleanupOldFiles()
    }

    lr.sequence++
    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", lr.basePath, lr.sequence)
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
    oldName := fmt.Sprintf("%s.%d.log", lr.basePath, lr.sequence-1)
    newName := oldName + ".gz"

    src, err := os.Open(oldName)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(newName)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(oldName)
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := fmt.Sprintf("%s.*.log.gz", lr.basePath)
    files, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(files) > maxBackups {
        filesToDelete := files[:len(files)-maxBackups]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
    }

    fmt.Println("Log rotation completed")
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	maxFileSize   = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logFileName   = "app.log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator(logDir string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fullPath := filepath.Join(logDir, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		basePath:    logDir,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	backupPath := filepath.Join(lr.basePath, backupName)

	if err := os.Rename(filepath.Join(lr.basePath, logFileName), backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	file, err := os.OpenFile(filepath.Join(lr.basePath, logFileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	lr.currentFile = file
	lr.currentSize = 0

	go lr.cleanupOldBackups()

	return nil
}

func (lr *LogRotator) cleanupOldBackups() {
	pattern := filepath.Join(lr.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackupFiles {
		return
	}

	sort.Strings(matches)
	filesToRemove := matches[:len(matches)-maxBackupFiles]

	for _, file := range filesToRemove {
		os.Remove(file)
	}
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("./logs")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 1024 * 1024 // 1MB
	maxBackups  = 5
)

type RotatingLogger struct {
	mu       sync.Mutex
	file     *os.File
	size     int64
	basePath string
	current  string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.current = rl.basePath + ".log"
	file, err := os.OpenFile(rl.current, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	rotated := fmt.Sprintf("%s.%s", rl.basePath, timestamp)
	if err := os.Rename(rl.current, rotated); err != nil {
		return err
	}

	if err := rl.cleanOldBackups(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clean old backups: %v\n", err)
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := rl.basePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackups {
		return nil
	}

	toRemove := matches[:len(matches)-maxBackups]
	for _, path := range toRemove {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
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
	currentFile *os.File
	basePath    string
	maxSize     int64
	currentSize int64
	fileCount   int
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
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rl.currentFile = f
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

	archiveName := fmt.Sprintf("%s.%d.%s.gz",
		rl.basePath,
		rl.fileCount,
		time.Now().Format("20060102T150405"))

	rl.fileCount++

	if err := compressFile(rl.basePath, archiveName); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrentFile()
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

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
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

type RotatingLog struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    rotationCount int
}

func NewRotatingLog(filePath string, maxSizeMB int) (*RotatingLog, error) {
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

    return &RotatingLog{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (r *RotatingLog) Write(p []byte) (int, error) {
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

func (r *RotatingLog) rotate() error {
    if r.file != nil {
        r.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", r.filePath, timestamp)

    oldFile, err := os.Open(r.filePath)
    if err != nil {
        return err
    }
    defer oldFile.Close()

    archivedFile, err := os.Create(archivedPath)
    if err != nil {
        return err
    }
    defer archivedFile.Close()

    gzWriter := gzip.NewWriter(archivedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }

    os.Remove(r.filePath)

    newFile, err := os.Create(r.filePath)
    if err != nil {
        return err
    }

    r.file = newFile
    r.currentSize = 0
    r.rotationCount++

    r.cleanOldArchives(5)

    return nil
}

func (r *RotatingLog) cleanOldArchives(keep int) {
    pattern := r.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= keep {
        return
    }

    for i := 0; i < len(matches)-keep; i++ {
        os.Remove(matches[i])
    }
}

func (r *RotatingLog) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: Application event occurred at %s\n", i, time.Now().Format(time.RFC3339))
        log.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}