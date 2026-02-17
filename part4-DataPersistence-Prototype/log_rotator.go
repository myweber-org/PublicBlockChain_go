
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
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
	currentSize int64
}

func NewRotatingLogger(filePath string, maxSizeMB int, backupCount int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	for i := rl.backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == rl.backupCount-1 {
				os.Remove(oldPath)
			} else {
				if err := rl.compressFile(oldPath, newPath); err != nil {
					return err
				}
			}
		}
	}

	if err := os.Rename(rl.filePath, rl.getBackupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.filePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.filePath, index+1)
}

func (rl *RotatingLogger) compressFile(src, dst string) error {
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

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	if err != nil {
		return err
	}

	os.Remove(src)
	return nil
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(message))
		time.Sleep(100 * time.Millisecond)
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

func (r *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(r.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    f, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    
    r.currentFile = f
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    
    if err := os.Rename(r.basePath, rotatedPath); err != nil {
        return err
    }
    
    if r.compressOld {
        go r.compressFile(rotatedPath)
    }
    
    if err := r.cleanupOldBackups(); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }
    
    return r.openCurrentFile()
}

func (r *LogRotator) compressFile(path string) {
    compressedPath := path + ".gz"
    
    src, err := os.Open(path)
    if err != nil {
        return
    }
    defer src.Close()
    
    dst, err := os.Create(compressedPath)
    if err != nil {
        return
    }
    defer dst.Close()
    
    gz := gzip.NewWriter(dst)
    defer gz.Close()
    
    if _, err := io.Copy(gz, src); err != nil {
        return
    }
    
    os.Remove(path)
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
    
    for i := r.maxBackups; i < len(backups); i++ {
        path := filepath.Join(dir, backups[i])
        os.Remove(path)
        if r.compressOld {
            os.Remove(path + ".gz")
        }
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

func extractTimestamp(name string) int64 {
    parts := strings.Split(name, ".")
    if len(parts) < 2 {
        return 0
    }
    
    timestamp := parts[len(parts)-1]
    if len(timestamp) != 14 {
        return 0
    }
    
    ts, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return 0
    }
    
    return ts
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
        logLine := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logLine))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation example completed")
}