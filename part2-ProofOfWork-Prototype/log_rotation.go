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
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type RotatingLogger struct {
    basePath      string
    maxSize       int64
    currentSize   int64
    currentFile   *os.File
    fileIndex     int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }
    
    if err := rl.openNewFile(); err != nil {
        return nil, err
    }
    
    return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    filename := fmt.Sprintf("%s_%d_%s.log", 
        rl.basePath, 
        rl.fileIndex, 
        time.Now().Format("20060102_150405"))
    
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    rl.currentFile = file
    rl.currentSize = 0
    rl.fileIndex++
    
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.openNewFile(); err != nil {
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
    logger, err := NewRotatingLogger("app_log", 1024*1024) // 1MB max size
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample data here\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}
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
	maxBackups  = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentIdx int
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
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	// Compress current file
	srcPath := rl.basePath
	dstPath := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentIdx)

	if err := compressFile(srcPath, dstPath); err != nil {
		return err
	}

	// Remove old backups
	rl.currentIdx = (rl.currentIdx + 1) % maxBackups
	oldPath := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentIdx)
	if _, err := os.Stat(oldPath); err == nil {
		os.Remove(oldPath)
	}

	// Remove uncompressed source
	os.Remove(srcPath)

	// Create new log file
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

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Simulate log writing
	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed. Check compressed files:")
	matches, _ := filepath.Glob("app.log.*.gz")
	for _, match := range matches {
		fmt.Println("  ", match)
	}
}
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
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu         sync.Mutex
	current    *os.File
	size       int64
	baseName   string
	sequence   int
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: filepath.Join(logDir, name),
	}

	if err := rl.openNew(); err != nil {
		return nil, err
	}

	return rl, nil
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
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.current != nil {
		rl.current.Close()
		rl.compressCurrent()
	}

	rl.sequence++
	if rl.sequence > maxBackups {
		rl.cleanOld()
	}

	return rl.openNew()
}

func (rl *RotatingLogger) openNew() error {
	filename := fmt.Sprintf("%s.%d.log", rl.baseName, rl.sequence)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.current = f
	if info, err := f.Stat(); err == nil {
		rl.size = info.Size()
	} else {
		rl.size = 0
	}

	return nil
}

func (rl *RotatingLogger) compressCurrent() {
	oldName := fmt.Sprintf("%s.%d.log", rl.baseName, rl.sequence)
	newName := oldName + ".gz"

	oldFile, err := os.Open(oldName)
	if err != nil {
		log.Printf("Failed to open for compression: %v", err)
		return
	}
	defer oldFile.Close()

	newFile, err := os.Create(newName)
	if err != nil {
		log.Printf("Failed to create compressed file: %v", err)
		return
	}
	defer newFile.Close()

	gz := gzip.NewWriter(newFile)
	defer gz.Close()

	if _, err := io.Copy(gz, oldFile); err != nil {
		log.Printf("Compression failed: %v", err)
		return
	}

	os.Remove(oldName)
}

func (rl *RotatingLogger) cleanOld() {
	for i := rl.sequence - maxBackups; i > 0; i-- {
		pattern := fmt.Sprintf("%s.%d.log*", rl.baseName, i)
		matches, _ := filepath.Glob(pattern)
		for _, match := range matches {
			os.Remove(match)
		}
	}
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}
}