
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

type LogRotator struct {
    mu          sync.Mutex
    currentFile *os.File
    basePath    string
    maxSize     int64
    fileCount   int
    maxFiles    int
}

func NewLogRotator(basePath string, maxSize int64, maxFiles int) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    fileInfo, err := lr.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if fileInfo.Size()+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    return lr.currentFile.Write(p)
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", lr.basePath, timestamp)

    if err := compressFile(lr.basePath, archivedPath); err != nil {
        return err
    }

    if err := os.Remove(lr.basePath); err != nil && !os.IsNotExist(err) {
        return err
    }

    lr.fileCount++
    if lr.fileCount > lr.maxFiles {
        if err := lr.cleanupOldFiles(); err != nil {
            return err
        }
    }

    return lr.openCurrentFile()
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

func (lr *LogRotator) cleanupOldFiles() error {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) > lr.maxFiles {
        filesToRemove := matches[:len(matches)-lr.maxFiles]
        for _, file := range filesToRemove {
            if err := os.Remove(file); err != nil {
                return err
            }
        }
    }

    return nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    lr.currentFile = file
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
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10*1024*1024, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}

	if err := rotator.openFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openFile() error {
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.file = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
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
	archivePath := fmt.Sprintf("%s.%s.gz", lr.filePath, timestamp)

	if err := compressFile(lr.filePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(lr.filePath); err != nil {
		return err
	}

	return lr.openFile()
}

func compressFile(source, target string) error {
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

func (lr *LogRotator) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}