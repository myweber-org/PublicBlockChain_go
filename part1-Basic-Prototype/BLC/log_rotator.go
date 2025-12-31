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

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		filePath: basePath,
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

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	rl.rotationCount++
	archiveName := fmt.Sprintf("%s.%d.%s.gz", 
		rl.filePath, 
		rl.rotationCount, 
		time.Now().Format("20060102_150405"))

	if err := rl.compressFile(rl.filePath, archiveName); err != nil {
		return err
	}

	if err := os.Truncate(rl.filePath, 0); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
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
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
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
	"strconv"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	backupCount int
	currentFile *os.File
	currentSize int64
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) (*LogRotator, error) {
	lr := &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := lr.openCurrentFile(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	for i := lr.backupCount - 1; i >= 0; i-- {
		oldPath := lr.getBackupPath(i)
		newPath := lr.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(lr.filePath, lr.getBackupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) getBackupPath(index int) string {
	if index == 0 {
		return lr.filePath + ".1"
	}
	return lr.filePath + "." + strconv.Itoa(index+1)
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
		lr.currentSize = 0
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 1024, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message for testing rotation.\n", i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
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
	rl.size += int64(n)
	return n, err
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

		filename := fmt.Sprintf("%s.%s.log", rl.basePath, dateStr)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

	_, err = io.Copy(gz, src)
	if err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	pattern := rl.basePath + ".*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, file := range toDelete {
			os.Remove(file)
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
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
)

type LogRotator struct {
    filePath string
    file     *os.File
    size     int64
}

func NewLogRotator(path string) (*LogRotator, error) {
    rotator := &LogRotator{filePath: path}
    if err := rotator.openFile(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.size+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := lr.file.Write(p)
    if err == nil {
        lr.size += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.file != nil {
        lr.file.Close()
    }

    baseDir := filepath.Dir(lr.filePath)
    baseName := filepath.Base(lr.filePath)
    nameWithoutExt := strings.TrimSuffix(baseName, logExtension)

    files, err := filepath.Glob(filepath.Join(baseDir, nameWithoutExt+".*"+logExtension))
    if err != nil {
        return err
    }

    var backups []string
    for _, f := range files {
        if strings.HasPrefix(filepath.Base(f), nameWithoutExt+".") {
            backups = append(backups, f)
        }
    }

    sort.Sort(sort.Reverse(sort.StringSlice(backups)))

    for i, backup := range backups {
        if i >= maxBackups-1 {
            os.Remove(backup)
            continue
        }
        parts := strings.Split(filepath.Base(backup), ".")
        if len(parts) < 3 {
            continue
        }
        num, err := strconv.Atoi(parts[1])
        if err != nil {
            continue
        }
        newName := filepath.Join(baseDir, fmt.Sprintf("%s.%d%s", nameWithoutExt, num+1, logExtension))
        os.Rename(backup, newName)
    }

    newBackup := filepath.Join(baseDir, fmt.Sprintf("%s.1%s", nameWithoutExt, logExtension))
    if err := os.Rename(lr.filePath, newBackup); err != nil && !os.IsNotExist(err) {
        return err
    }

    return lr.openFile()
}

func (lr *LogRotator) openFile() error {
    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    lr.file = file
    lr.size = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format("2006-01-02 15:04:05"), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}