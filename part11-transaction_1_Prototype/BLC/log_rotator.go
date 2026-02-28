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
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath:  basePath,
        maxSize:   maxSize,
        fileIndex: 0,
    }

    err := logger.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileIndex)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

    oldFilename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileIndex)
    compressedFilename := fmt.Sprintf("%s.%d.log.gz", rl.basePath, rl.fileIndex)

    err := compressFile(oldFilename, compressedFilename)
    if err != nil {
        return err
    }

    err = os.Remove(oldFilename)
    if err != nil {
        return err
    }

    rl.fileIndex++
    return rl.openCurrentFile()
}

func compressFile(source, target string) error {
    sourceFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    targetFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer targetFile.Close()

    gzWriter := gzip.NewWriter(targetFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
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
    logger, err := NewRotatingLogger("app", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: This is a sample log message.\n",
            time.Now().Format(time.RFC3339), i)
        _, err := logger.Write([]byte(message))
        if err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed. Check generated files.")
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
    fileIndex   int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
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
        if err := lr.compressFile(lr.currentFile.Name()); err != nil {
            return err
        }
    }

    lr.fileIndex++
    if lr.fileIndex > maxBackups {
        lr.fileIndex = 1
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", lr.basePath, lr.fileIndex)
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

func (lr *LogRotator) compressFile(source string) error {
    dest := source + ".gz"
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    return os.Remove(source)
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.basePath + ".*.log.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        stat, err := os.Stat(match)
        if err != nil {
            continue
        }
        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{match, stat.ModTime()})
    }

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        oldest := backupFiles[i]
        os.Remove(oldest.path)
    }

    return nil
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
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    rotator.cleanupOldBackups()
    fmt.Println("Log rotation completed successfully")
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
	maxFileSize = 10 * 1024 * 1024
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
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

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.size = stat.Size()
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

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == backupCount-1 {
				os.Remove(oldPath)
			} else {
				if err := rl.compressFile(oldPath, newPath); err != nil {
					return err
				}
			}
		}
	}

	if err := os.Rename(rl.basePath, rl.getBackupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) getBackupPath(num int) string {
	if num == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.basePath, num+1)
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

	if err := srcFile.Close(); err != nil {
		return err
	}

	return os.Remove(src)
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
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Application event recorded\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(msg))
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
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentSize int64
	basePath    string
	file        *os.File
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
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
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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
		rl.file = nil
	}

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == backupCount-1 {
				os.Remove(oldPath)
			} else {
				if err := rl.compressFile(oldPath, newPath); err != nil {
					return err
				}
			}
		}
	}

	if err := os.Rename(rl.basePath, rl.getBackupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.basePath, index+1)
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

	return os.Remove(src)
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
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
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
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	currentSize int64
	basePath    string
	fileIndex   int
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
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
	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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

	oldPath := rl.basePath + ".log"
	newPath := fmt.Sprintf("%s.%d.log.gz", rl.basePath, rl.fileIndex)

	if err := compressFile(oldPath, newPath); err != nil {
		return err
	}

	if err := os.Remove(oldPath); err != nil {
		return err
	}

	rl.fileIndex = (rl.fileIndex + 1) % maxBackups
	return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
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
	logger, err := NewRotatingLogger("application")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Sample log message for testing rotation\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed. Check generated files.")
}package main

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

type RotatingLog struct {
    filePath   string
    current    *os.File
    currentSize int64
}

func NewRotatingLog(path string) (*RotatingLog, error) {
    rl := &RotatingLog{filePath: path}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := rl.current.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLog) rotate() error {
    if err := rl.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    if err := os.Rename(rl.filePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := rl.cleanupOldBackups(); err != nil {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLog) compressFile(path string) error {
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

    if err := os.Remove(path); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLog) cleanupOldBackups() error {
    dir := filepath.Dir(rl.filePath)
    base := filepath.Base(rl.filePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, base+".") && strings.HasSuffix(name, ".gz") {
            backups = append(backups, name)
        }
    }

    if len(backups) <= maxBackups {
        return nil
    }

    sortBackups(backups)
    for i := maxBackups; i < len(backups); i++ {
        if err := os.Remove(filepath.Join(dir, backups[i])); err != nil {
            return err
        }
    }

    return nil
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            timeI := extractTimestamp(backups[i])
            timeJ := extractTimestamp(backups[j])
            if timeI < timeJ {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return 0
    }
    ts, _ := strconv.ParseInt(parts[len(parts)-2], 10, 64)
    return ts
}

func (rl *RotatingLog) openCurrent() error {
    file, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.current = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLog) Close() error {
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("application.log")
    if err != nil {
        fmt.Printf("Failed to create log: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type LogRotator struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
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
    
    return &LogRotator{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    
    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)
    
    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    lr.file = file
    lr.currentSize = 0
    
    go lr.cleanOldLogs()
    
    return nil
}

func (lr *LogRotator) cleanOldLogs() {
    pattern := fmt.Sprintf("%s.*", lr.filePath)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) > 10 {
        for i := 0; i < len(matches)-10; i++ {
            os.Remove(matches[i])
        }
    }
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}