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
		rl.currentFile = nil
	}

	rl.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.%d.gz", rl.filePath, timestamp, rl.rotationCount)

	if err := rl.compressFile(rl.filePath, archivePath); err != nil {
		return err
	}

	if err := os.Truncate(rl.filePath, 0); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
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
	logger, err := NewRotatingLogger("application.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Sample log data for testing rotation\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
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

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

	dst, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
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
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
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
    "sync"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    mu         sync.Mutex
    file       *os.File
    currentSize int64
    basePath   string
    sequence   int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
        sequence: 0,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > maxFileSize {
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
    if lr.file != nil {
        lr.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.%d", lr.basePath, timestamp, lr.sequence)
    
    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    lr.sequence++
    if lr.sequence > maxBackups {
        lr.cleanupOldFiles()
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(source string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(source + ".gz")
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

    os.Remove(source)
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        sortFiles(matches)
        for i := 0; i < len(matches)-maxBackups; i++ {
            os.Remove(matches[i])
        }
    }
}

func sortFiles(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractSequence(files[i]) > extractSequence(files[j]) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractSequence(filename string) int {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return 0
    }
    seqStr := parts[len(parts)-2]
    seq, _ := strconv.Atoi(seqStr)
    return seq
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
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
	backupCount = 5
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
	currentDay := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDay {
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
		rl.currentDay = currentDay
	}
	return nil
}

func (rl *RotatingLogger) getLogPath(t time.Time) string {
	if rl.size >= maxFileSize || rl.file == nil {
		return fmt.Sprintf("%s.%s.%d.log", rl.basePath, t.Format("2006-01-02"), t.Unix())
	}
	return rl.basePath + ".log"
}

func (rl *RotatingLogger) compressOldLog() error {
	oldPath := rl.file.Name()
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

	if len(matches) > backupCount {
		toDelete := matches[:len(matches)-backupCount]
		for _, file := range toDelete {
			os.Remove(file)
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
	logger, err := NewRotatingLogger("/var/log/myapp/application")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 100; i++ {
		customLog.Printf("Log entry %d: Application is running smoothly", i)
		time.Sleep(100 * time.Millisecond)
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

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    sequence    int
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
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
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
    rl.sequence = 0
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
    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.file != nil {
        rl.file.Close()
    }
    timestamp := time.Now().Format("20060102_150405")
    archiveName := fmt.Sprintf("%s.%s.%d.gz", rl.basePath, timestamp, rl.sequence)
    if err := rl.compressFile(rl.basePath, archiveName); err != nil {
        return err
    }
    if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
        return err
    }
    rl.sequence++
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
    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("/var/log/app/application.log", 10)
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
}package main

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
    mu           sync.Mutex
    basePath     string
    maxSize      int64
    currentFile  *os.File
    currentSize  int64
    maxBackups   int
    compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxBackups int, compress bool) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compress,
    }

    err := logger.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = info.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (n int, err error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        err := l.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err = l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    backupPath := l.basePath + "." + timestamp

    err := os.Rename(l.basePath, backupPath)
    if err != nil {
        return err
    }

    err = l.openCurrentFile()
    if err != nil {
        return err
    }

    go l.manageBackups(backupPath)

    return nil
}

func (l *RotatingLogger) manageBackups(backupPath string) {
    if l.compressOld {
        compressedPath := backupPath + ".gz"
        err := compressFile(backupPath, compressedPath)
        if err == nil {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }

    backups, err := l.listBackups()
    if err != nil {
        return
    }

    if len(backups) > l.maxBackups {
        backupsToDelete := backups[:len(backups)-l.maxBackups]
        for _, backup := range backupsToDelete {
            os.Remove(backup)
        }
    }
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

func (l *RotatingLogger) listBackups() ([]string, error) {
    pattern := l.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return nil, err
    }

    var backups []string
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") || isTimestampBackup(match, l.basePath) {
            backups = append(backups, match)
        }
    }

    sortBackupsByTimestamp(backups, l.basePath)
    return backups, nil
}

func isTimestampBackup(path, basePath string) bool {
    suffix := strings.TrimPrefix(path, basePath+".")
    if len(suffix) != 14 {
        return false
    }
    _, err := strconv.ParseInt(suffix, 10, 64)
    return err == nil
}

func sortBackupsByTimestamp(backups []string, basePath string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            timeI := extractTimestamp(backups[i], basePath)
            timeJ := extractTimestamp(backups[j], basePath)
            if timeI < timeJ {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(path, basePath string) int64 {
    suffix := strings.TrimSuffix(strings.TrimPrefix(path, basePath+"."), ".gz")
    timestamp, _ := strconv.ParseInt(suffix, 10, 64)
    return timestamp
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(message))
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
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingWriter struct {
    filename   string
    current    *os.File
    size       int64
    mu         sync.Mutex
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
    w := &RotatingWriter{filename: filename}
    if err := w.openFile(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.size+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.current.Write(p)
    w.size += int64(n)
    return n, err
}

func (w *RotatingWriter) openFile() error {
    file, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
    return nil
}

func (w *RotatingWriter) rotate() error {
    if w.current != nil {
        w.current.Close()
    }

    timestamp := time.Now().Format("20060102-150405")
    backupName := fmt.Sprintf("%s.%s.gz", w.filename, timestamp)

    if err := compressFile(w.filename, backupName); err != nil {
        return err
    }

    if err := cleanupOldBackups(w.filename); err != nil {
        return err
    }

    os.Remove(w.filename)
    return w.openFile()
}

func compressFile(source, target string) error {
    src, err := os.Open(source)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    return err
}

func cleanupOldBackups(baseName string) error {
    pattern := fmt.Sprintf("%s.*.gz", filepath.Base(baseName))
    matches, err := filepath.Glob(filepath.Join(filepath.Dir(baseName), pattern))
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
    return nil
}

func (w *RotatingWriter) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    if w.current != nil {
        return w.current.Close()
    }
    return nil
}