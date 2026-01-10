
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

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: name,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	if info, err := file.Stat(); err == nil {
		rl.currentSize = info.Size()
	}
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return err
	}

	if len(files) >= maxBackups {
		oldest := files[0]
		for _, f := range files[1:] {
			if fi1, _ := os.Stat(f); fi2, _ := os.Stat(oldest); fi1.ModTime().Before(fi2.ModTime()) {
				oldest = f
			}
		}
		if err := compressAndRemove(oldest); err != nil {
			return err
		}
	}

	return rl.openCurrentFile()
}

func compressAndRemove(filename string) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(filename)
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}package main

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

type RotatingLogger struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    logger := &RotatingLogger{
        filePath: basePath,
    }
    
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return logger, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.bytesWritten+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := rl.currentFile.Write(p)
    rl.bytesWritten += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.filePath, timestamp)
    
    if err := compressFile(rl.filePath, archivedPath); err != nil {
        return err
    }
    
    if err := cleanupOldBackups(rl.filePath); err != nil {
        return err
    }
    
    os.Remove(rl.filePath)
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
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
    rl.bytesWritten = info.Size()
    return nil
}

func compressFile(source, destination string) error {
    src, err := os.Open(source)
    if err != nil {
        return err
    }
    defer src.Close()
    
    dst, err := os.Create(destination)
    if err != nil {
        return err
    }
    defer dst.Close()
    
    gz := gzip.NewWriter(dst)
    defer gz.Close()
    
    _, err = io.Copy(gz, src)
    return err
}

func cleanupOldBackups(basePath string) error {
    pattern := basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    
    if len(matches) > maxBackups {
        filesToRemove := matches[:len(matches)-maxBackups]
        for _, file := range filesToRemove {
            os.Remove(file)
        }
    }
    
    return nil
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("application.log")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Sample log message\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation demonstration completed")
}package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize    = 10 * 1024 * 1024 // 10MB
	maxLogFiles   = 5
	logFileName   = "app.log"
	checkInterval = 30 * time.Second
)

func rotateLogs() error {
	info, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if info.Size() < maxLogSize {
		return nil
	}

	for i := maxLogFiles - 1; i > 0; i-- {
		oldName := logFileName + "." + string(rune('0'+i))
		newName := logFileName + "." + string(rune('0'+i+1))
		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	backupName := logFileName + ".1"
	os.Rename(logFileName, backupName)

	files, _ := filepath.Glob(logFileName + ".*")
	if len(files) > maxLogFiles {
		for _, f := range files[maxLogFiles:] {
			os.Remove(f)
		}
	}

	return nil
}

func main() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := rotateLogs(); err != nil {
			log.Printf("Log rotation failed: %v", err)
		}
	}
}