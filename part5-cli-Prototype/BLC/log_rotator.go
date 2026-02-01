
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
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
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

	path := rl.basePath + ".log"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
	if rl.file != nil {
		rl.file.Close()
	}

	oldPath := rl.basePath + ".log"
	newPath := fmt.Sprintf("%s.%d.log.gz", rl.basePath, time.Now().Unix())

	if err := compressFile(oldPath, newPath); err != nil {
		return err
	}

	if err := os.Truncate(oldPath, 0); err != nil {
		return err
	}

	rl.cleanOldBackups()
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

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := rl.basePath + ".*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		for i := 0; i < len(matches)-maxBackups; i++ {
			os.Remove(matches[i])
		}
	}
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(msg))
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
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile   *os.File
    currentSize   int64
    basePath      string
    currentNumber int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    rotator.findLatestRotationNumber()
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        err := lr.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    lr.currentSize += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    rotatedPath := lr.getRotatedPath()
    err := os.Rename(lr.basePath, rotatedPath)
    if err != nil {
        return err
    }

    err = lr.compressFile(rotatedPath)
    if err != nil {
        return err
    }

    lr.currentNumber++
    err = lr.cleanupOldBackups()
    if err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) getRotatedPath() string {
    timestamp := time.Now().Format("20060102_150405")
    return fmt.Sprintf("%s.%s.%d", lr.basePath, timestamp, lr.currentNumber)
}

func (lr *LogRotator) compressFile(sourcePath string) error {
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

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    backupNumbers := make([]int, 0, len(matches))
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        numStr := parts[len(parts)-2]
        num, err := strconv.Atoi(numStr)
        if err != nil {
            continue
        }
        backupNumbers = append(backupNumbers, num)
    }

    for i := 0; i < len(backupNumbers)-maxBackups; i++ {
        oldestNum := backupNumbers[i]
        pattern := fmt.Sprintf("%s.*.%d.gz", lr.basePath, oldestNum)
        matches, _ := filepath.Glob(pattern)
        for _, match := range matches {
            os.Remove(match)
        }
    }

    return nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (lr *LogRotator) findLatestRotationNumber() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil || len(matches) == 0 {
        lr.currentNumber = 1
        return
    }

    maxNum := 0
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        numStr := parts[len(parts)-2]
        num, err := strconv.Atoi(numStr)
        if err == nil && num > maxNum {
            maxNum = num
        }
    }
    lr.currentNumber = maxNum + 1
}

func (lr *LogRotator) Close() error {
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
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}