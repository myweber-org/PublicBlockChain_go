
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
    basePath     string
    maxSize      int64
    currentSize  int64
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%d.%s", rl.basePath, rl.rotationCount, timestamp)
    rl.rotationCount++

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    destFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
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
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Sample log message for testing rotation\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
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
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewLogRotator(name string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rotator := &LogRotator{
		baseName: name,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", lr.baseName, timestamp))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	if err := lr.compressOldLogs(); err != nil {
		return err
	}

	if err := lr.cleanupOldBackups(); err != nil {
		return err
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) compressOldLogs() error {
	files, err := filepath.Glob(filepath.Join(logDir, fmt.Sprintf("%s_*.log", lr.baseName)))
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file) == ".gz" {
			continue
		}

		compressedFile := file + ".gz"
		if _, err := os.Stat(compressedFile); err == nil {
			continue
		}

		if err := compressFile(file, compressedFile); err != nil {
			return err
		}

		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
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

func (lr *LogRotator) cleanupOldBackups() error {
	files, err := filepath.Glob(filepath.Join(logDir, fmt.Sprintf("%s_*.log.gz", lr.baseName)))
	if err != nil {
		return err
	}

	if len(files) <= maxBackups {
		return nil
	}

	for i := 0; i < len(files)-maxBackups; i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}

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
	rotator, err := NewLogRotator("application")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
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
    if maxSizeMB <= 0 {
        maxSizeMB = 10
    }
    if maxFiles <= 0 {
        maxFiles = 5
    }

    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  int64(maxSizeMB) * 1024 * 1024,
        maxFiles: maxFiles,
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

    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    stat, err := rl.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if stat.Size()+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    rl.fileCount++
    rl.cleanupOldFiles()

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    if err := os.Remove(sourcePath); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

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

    if len(compressedFiles) > rl.maxFiles {
        filesToRemove := len(compressedFiles) - rl.maxFiles
        for i := 0; i < filesToRemove; i++ {
            os.Remove(compressedFiles[i])
        }
    }
}

func (rl *RotatingLogger) parseFileNumber(filename string) int {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return 0
    }

    numStr := parts[len(parts)-2]
    num, err := strconv.Atoi(numStr)
    if err != nil {
        return 0
    }

    return num
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
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 5, 3)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
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
    "sync"
)

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  int64(maxSizeMB) * 1024 * 1024,
        sequence: 0,
    }

    if err := rl.openCurrent(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) openCurrent() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        rl.file.Close()
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

    rl.file = file
    rl.currentSize = info.Size()

    rl.findLatestSequence()
    return nil
}

func (rl *RotatingLog) findLatestSequence() {
    pattern := rl.basePath + ".*.gz"
    matches, _ := filepath.Glob(pattern)
    maxSeq := 0

    for _, match := range matches {
        base := filepath.Base(match)
        if len(base) > len(rl.basePath)+5 {
            seqStr := base[len(rl.basePath)+1 : len(base)-3]
            if seq, err := strconv.Atoi(seqStr); err == nil && seq > maxSeq {
                maxSeq = seq
            }
        }
    }

    rl.sequence = maxSeq
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
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

func (rl *RotatingLog) rotate() error {
    if rl.file != nil {
        rl.file.Close()
    }

    rl.sequence++
    archivedName := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.sequence)

    oldFile, err := os.Open(rl.basePath)
    if err != nil {
        return err
    }
    defer oldFile.Close()

    archivedFile, err := os.Create(archivedName)
    if err != nil {
        return err
    }
    defer archivedFile.Close()

    gzWriter := gzip.NewWriter(archivedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }

    os.Remove(rl.basePath)

    return rl.openCurrent()
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}