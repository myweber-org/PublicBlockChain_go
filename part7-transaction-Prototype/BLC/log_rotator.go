package main

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
    currentSize int64
    basePath    string
    sequence    int
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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
    
    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)
    
    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }
    
    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }
    
    rl.cleanOldBackups()
    
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source string) error {
    src, err := os.Open(source)
    if err != nil {
        return err
    }
    defer src.Close()
    
    dest, err := os.Create(source + ".gz")
    if err != nil {
        return err
    }
    defer dest.Close()
    
    gz := gzip.NewWriter(dest)
    defer gz.Close()
    
    _, err = io.Copy(gz, src)
    if err != nil {
        return err
    }
    
    os.Remove(source)
    return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) > maxBackups {
        toRemove := matches[:len(matches)-maxBackups]
        for _, file := range toRemove {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
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
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
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
	rl.currentSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
		timestamp := time.Now().Format("20060102_150405")
		oldPath := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))
		if err := os.Rename(rl.currentFile.Name(), oldPath); err != nil {
			return err
		}
		if err := rl.compressFile(oldPath); err != nil {
			return err
		}
		rl.cleanupOld()
	}
	return rl.openCurrent()
}

func (rl *RotatingLogger) openCurrent() error {
	path := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) compressFile(src string) error {
	dest := src + ".gz"
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gz := gzip.NewWriter(destFile)
	defer gz.Close()

	if _, err := io.Copy(gz, srcFile); err != nil {
		return err
	}
	os.Remove(src)
	return nil
}

func (rl *RotatingLogger) cleanupOld() {
	pattern := filepath.Join(logDir, rl.baseName+"_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}
	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, f := range toDelete {
			os.Remove(f)
		}
	}
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
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Log rotation test completed")
}