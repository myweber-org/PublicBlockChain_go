
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
)

type RotatingWriter struct {
    filename   string
    current    *os.File
    size       int64
    backupTime time.Time
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
    writer := &RotatingWriter{
        filename: filename,
    }
    
    if err := writer.openFile(); err != nil {
        return nil, err
    }
    
    return writer, nil
}

func (w *RotatingWriter) openFile() error {
    file, err := os.OpenFile(w.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    w.current = file
    w.size = info.Size()
    w.backupTime = time.Now()
    
    return nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.size+int64(len(p)) >= maxFileSize || time.Since(w.backupTime).Hours() >= 24 {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := w.current.Write(p)
    if err != nil {
        return n, err
    }
    
    w.size += int64(n)
    return n, nil
}

func (w *RotatingWriter) rotate() error {
    if w.current != nil {
        if err := w.current.Close(); err != nil {
            return err
        }
    }
    
    for i := maxBackups - 1; i >= 0; i-- {
        oldName := w.backupName(i)
        newName := w.backupName(i + 1)
        
        if _, err := os.Stat(oldName); err == nil {
            if err := os.Rename(oldName, newName); err != nil {
                return err
            }
        }
    }
    
    if err := os.Rename(w.filename, w.backupName(0)); err != nil && !os.IsNotExist(err) {
        return err
    }
    
    return w.openFile()
}

func (w *RotatingWriter) backupName(index int) string {
    if index == 0 {
        return w.filename + ".1"
    }
    return fmt.Sprintf("%s.%d", w.filename, index+1)
}

func (w *RotatingWriter) Close() error {
    if w.current != nil {
        return w.current.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log")
    if err != nil {
        log.Fatal(err)
    }
    defer writer.Close()
    
    logger := log.New(writer, "", log.LstdFlags)
    
    for i := 0; i < 1000; i++ {
        logger.Printf("Log entry %d: %s", i, "Sample log message for testing rotation")
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if err := rl.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err = rl.file.Write(p)
    rl.size += int64(n)
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
            rl.cleanupOldBackups()
        }

        newPath := rl.getLogPath(now)
        if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
            return err
        }

        file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
        rl.currentDay = today
    }
    return nil
}

func (rl *RotatingLogger) getLogPath(t time.Time) string {
    if rl.size >= maxFileSize {
        timestamp := t.Format("20060102-150405")
        return fmt.Sprintf("%s.%s.log", rl.basePath, timestamp)
    }
    return rl.basePath + ".log"
}

func (rl *RotatingLogger) compressOldLog() error {
    oldPath := rl.getLogPath(time.Now().Add(-time.Second))
    if _, err := os.Stat(oldPath); os.IsNotExist(err) {
        return nil
    }

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
        for _, f := range toDelete {
            os.Remove(f)
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
    logger, err := NewRotatingLogger("/var/log/myapp/app")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(logger)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(100 * time.Millisecond)
    }
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
}

func NewRotatingLogger(path string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	
	return &RotatingLogger{
		filePath:    path,
		maxSize:     maxSize,
		currentSize: info.Size(),
		file:        file,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()
	
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, time.Now().Format("20060102150405"))
	
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}
	
	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	rl.file = file
	rl.currentSize = 0
	
	rl.cleanOldBackups()
	return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
	dir := filepath.Dir(rl.filePath)
	base := filepath.Base(rl.filePath)
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	
	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		if len(name) > len(base)+15 && name[:len(base)] == base && name[len(base)] == '.' {
			backups = append(backups, name)
		}
	}
	
	if len(backups) > 5 {
		oldest := backups[0]
		os.Remove(filepath.Join(dir, oldest))
	}
}

func (rl *RotatingLogger) Close() error {
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()
	
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application event occurred at %s\n", 
			i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}