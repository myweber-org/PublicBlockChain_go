
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "log"
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
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
    now := time.Now()
    today := now.Format("2006-01-02")

    if rl.file == nil || rl.currentDay != today || rl.size >= maxFileSize {
        if rl.file != nil {
            rl.file.Close()
            if err := rl.compressOldLog(); err != nil {
                log.Printf("Failed to compress log: %v", err)
            }
        }

        newPath := fmt.Sprintf("%s.%s.log", rl.basePath, today)
        f, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            return err
        }

        info, err := f.Stat()
        if err != nil {
            f.Close()
            return err
        }

        rl.file = f
        rl.size = info.Size()
        rl.currentDay = today
        rl.cleanupOldBackups()
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

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.basePath + ".*.log.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, path := range toDelete {
            os.Remove(path)
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
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(logger)

    for i := 0; i < 100; i++ {
        log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(100 * time.Millisecond)
    }
}package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const maxLogSize = 1024 * 1024 // 1MB
const backupCount = 5

type RotatingLogger struct {
	filePath string
	currentSize int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{filePath: path}
	
	info, err := os.Stat(path)
	if err == nil {
		rl.currentSize = info.Size()
	} else if os.IsNotExist(err) {
		rl.currentSize = 0
	} else {
		return nil, err
	}
	
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxLogSize {
		err := rl.rotate()
		if err != nil {
			return 0, err
		}
	}
	
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	
	n, err := file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	// Close current file
	rl.currentSize = 0
	
	// Rotate backup files
	for i := backupCount - 1; i >= 0; i-- {
		var source, dest string
		
		if i == 0 {
			source = rl.filePath
		} else {
			source = rl.filePath + "." + strconv.Itoa(i)
		}
		
		dest = rl.filePath + "." + strconv.Itoa(i+1)
		
		if _, err := os.Stat(source); err == nil {
			err := os.Rename(source, dest)
			if err != nil {
				return err
			}
		}
	}
	
	// Create new log file
	file, err := os.Create(rl.filePath)
	if err != nil {
		return err
	}
	file.Close()
	
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	
	log.SetOutput(logger)
	
	for i := 0; i < 1000; i++ {
		log.Printf("Log entry number %d", i)
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

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	maxBackups   int
	currentSize  int64
}

func NewRotatingLogger(filePath string, maxSize int64, maxBackups int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath:   filePath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	dir := filepath.Dir(rl.filePath)
	base := filepath.Base(rl.filePath)

	for i := rl.maxBackups - 1; i >= 0; i-- {
		oldPath := filepath.Join(dir, fmt.Sprintf("%s.%d", base, i))
		newPath := filepath.Join(dir, fmt.Sprintf("%s.%d", base, i+1))

		if i == rl.maxBackups-1 {
			os.Remove(oldPath)
		} else {
			if _, err := os.Stat(oldPath); err == nil {
				os.Rename(oldPath, newPath)
			}
		}
	}

	if _, err := os.Stat(rl.filePath); err == nil {
		backupPath := filepath.Join(dir, fmt.Sprintf("%s.0", base))
		os.Rename(rl.filePath, backupPath)
	}

	return rl.openCurrentFile()
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
	logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}