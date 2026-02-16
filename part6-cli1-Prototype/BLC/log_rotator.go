
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

const (
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    mu          sync.Mutex
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    lr := &LogRotator{
        basePath: basePath,
    }

    if err := lr.openCurrentFile(); err != nil {
        return nil, err
    }

    return lr, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

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
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, backupPath); err != nil {
        return err
    }

    if err := lr.compressBackup(backupPath); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressBackup(backupPath string) error {
    srcFile, err := os.Open(backupPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    gzPath := backupPath + ".gz"
    dstFile, err := os.Create(gzPath)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    os.Remove(backupPath)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackupFiles {
        return nil
    }

    var timestamps []time.Time
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        tsStr := parts[len(parts)-2]
        t, err := time.Parse("20060102_150405", tsStr)
        if err != nil {
            continue
        }
        timestamps = append(timestamps, t)
    }

    for i := 0; i < len(timestamps)-maxBackupFiles; i++ {
        oldestIdx := 0
        for j := 1; j < len(timestamps); j++ {
            if timestamps[j].Before(timestamps[oldestIdx]) {
                oldestIdx = j
            }
        }
        os.Remove(matches[oldestIdx])
        timestamps = append(timestamps[:oldestIdx], timestamps[oldestIdx+1:]...)
        matches = append(matches[:oldestIdx], matches[oldestIdx+1:]...)
    }

    return nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

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
            fmt.Printf("Write error: %v\n", err)
            break
        }

        if i%100 == 0 {
            fmt.Printf("Written %d log entries\n", i+1)
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
	"sync"
)

type RotatingWriter struct {
	mu          sync.Mutex
	file        *os.File
	maxSize     int64
	basePath    string
	currentSize int64
	maxFiles    int
}

func NewRotatingWriter(basePath string, maxSize int64, maxFiles int) (*RotatingWriter, error) {
	writer := &RotatingWriter{
		basePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := writer.openCurrentFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *RotatingWriter) openCurrentFile() error {
	file, err := os.OpenFile(w.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.file = file
	w.currentSize = stat.Size()
	return nil
}

func (w *RotatingWriter) rotate() error {
	w.file.Close()

	for i := w.maxFiles - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", w.basePath, i)
		newPath := fmt.Sprintf("%s.%d", w.basePath, i+1)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	backupPath := fmt.Sprintf("%s.1", w.basePath)
	os.Rename(w.basePath, backupPath)

	return w.openCurrentFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}

func main() {
	writer, err := NewRotatingWriter("app.log", 1024*1024, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create rotating writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		line := fmt.Sprintf("Log entry number %d: Some sample log data for testing rotation\n", i)
		if _, err := writer.Write([]byte(line)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			break
		}
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
    "sync"
    "time"
)

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compressOld bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }
    
    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := r.currentFile.Write(p)
    if err != nil {
        return n, err
    }
    
    r.currentSize += int64(n)
    return n, nil
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    
    if err := os.Rename(r.basePath, rotatedPath); err != nil {
        return err
    }
    
    if r.compressOld {
        go r.compressFile(rotatedPath)
    }
    
    if err := r.cleanupOldBackups(); err != nil {
        return err
    }
    
    return r.openCurrentFile()
}

func (r *LogRotator) compressFile(path string) error {
    srcFile, err := os.Open(path)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    
    dstFile, err := os.Create(path + ".gz")
    if err != nil {
        return err
    }
    defer dstFile.Close()
    
    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()
    
    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }
    
    os.Remove(path)
    return nil
}

func (r *LogRotator) cleanupOldBackups() error {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)
    
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }
    
    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            backups = append(backups, name)
        }
    }
    
    if len(backups) <= r.maxBackups {
        return nil
    }
    
    sortBackups(backups)
    
    for i := 0; i < len(backups)-r.maxBackups; i++ {
        os.Remove(filepath.Join(dir, backups[i]))
    }
    
    return nil
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return 0
    }
    
    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }
    
    timestamp, err := strconv.ParseInt(strings.ReplaceAll(timestampStr, "_", ""), 10, 64)
    if err != nil {
        return 0
    }
    
    return timestamp
}

func (r *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5, true)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}