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
	mu            sync.Mutex
	currentFile   *os.File
	basePath      string
	maxSize       int64
	currentSize   int64
	rotationCount int
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

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
	archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	if err := rl.compressFile(rl.basePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	rl.rotationCount++
	return rl.openCurrentFile()
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

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    FilePath    string
    MaxSize     int64
    MaxFiles    int
    RotateEvery time.Duration
    lastRotate  time.Time
}

func NewLogRotator(path string, maxSize int64, maxFiles int, rotateEvery time.Duration) *LogRotator {
    return &LogRotator{
        FilePath:    path,
        MaxSize:     maxSize,
        MaxFiles:    maxFiles,
        RotateEvery: rotateEvery,
        lastRotate:  time.Now(),
    }
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if err := lr.checkRotation(); err != nil {
        return 0, err
    }

    file, err := os.OpenFile(lr.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    return file.Write(p)
}

func (lr *LogRotator) checkRotation() error {
    now := time.Now()
    shouldRotate := false

    if lr.RotateEvery > 0 && now.Sub(lr.lastRotate) >= lr.RotateEvery {
        shouldRotate = true
        lr.lastRotate = now
    }

    if !shouldRotate && lr.MaxSize > 0 {
        if info, err := os.Stat(lr.FilePath); err == nil && info.Size() >= lr.MaxSize {
            shouldRotate = true
        }
    }

    if shouldRotate {
        return lr.performRotation()
    }
    return nil
}

func (lr *LogRotator) performRotation() error {
    for i := lr.MaxFiles - 1; i > 0; i-- {
        oldName := fmt.Sprintf("%s.%d", lr.FilePath, i)
        newName := fmt.Sprintf("%s.%d", lr.FilePath, i+1)

        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    if _, err := os.Stat(lr.FilePath); err == nil {
        return os.Rename(lr.FilePath, fmt.Sprintf("%s.1", lr.FilePath))
    }
    return nil
}

func (lr *LogRotator) Cleanup() error {
    for i := lr.MaxFiles + 1; ; i++ {
        fileName := fmt.Sprintf("%s.%d", lr.FilePath, i)
        if _, err := os.Stat(fileName); os.IsNotExist(err) {
            break
        }
        os.Remove(fileName)
    }
    return nil
}

func main() {
    rotator := NewLogRotator(
        "/var/log/app.log",
        10*1024*1024,
        5,
        time.Hour*24,
    )

    message := fmt.Sprintf("[%s] Application started\n", time.Now().Format(time.RFC3339))
    rotator.Write([]byte(message))

    fmt.Println("Log entry written successfully")
}