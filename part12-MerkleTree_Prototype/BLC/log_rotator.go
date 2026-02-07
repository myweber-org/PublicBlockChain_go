package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type RotatingLogger struct {
	dirPath      string
	baseName     string
	maxFileSize  int64
	maxFiles     int
	currentFile  *os.File
	currentSize  int64
}

func NewRotatingLogger(dir, base string, maxSize int64, maxCount int) (*RotatingLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	rl := &RotatingLogger{
		dirPath:     dir,
		baseName:    base,
		maxFileSize: maxSize,
		maxFiles:    maxCount,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	pattern := filepath.Join(rl.dirPath, rl.baseName+".log*")
	matches, _ := filepath.Glob(pattern)

	var maxNum int
	for _, m := range matches {
		var num int
		fmt.Sscanf(filepath.Ext(m), ".log%d", &num)
		if num > maxNum {
			maxNum = num
		}
	}

	filename := filepath.Join(rl.dirPath, fmt.Sprintf("%s.log", rl.baseName))
	if maxNum > 0 {
		filename = filepath.Join(rl.dirPath, fmt.Sprintf("%s.log%d", rl.baseName, maxNum))
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.currentFile = f

	info, err := f.Stat()
	if err != nil {
		return err
	}
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxFileSize {
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

	files, _ := filepath.Glob(filepath.Join(rl.dirPath, rl.baseName+".log*"))
	sort.Slice(files, func(i, j int) bool {
		return extractNum(files[i]) > extractNum(files[j])
	})

	for i, f := range files {
		oldNum := extractNum(f)
		if oldNum == 0 {
			newName := filepath.Join(rl.dirPath, fmt.Sprintf("%s.log1", rl.baseName))
			os.Rename(f, newName)
		} else if oldNum < rl.maxFiles {
			newName := filepath.Join(rl.dirPath, fmt.Sprintf("%s.log%d", rl.baseName, oldNum+1))
			os.Rename(f, newName)
		} else {
			os.Remove(f)
		}
		if i >= rl.maxFiles-1 {
			break
		}
	}

	return rl.openCurrent()
}

func extractNum(path string) int {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == ".log" {
		return 0
	}
	var num int
	fmt.Sscanf(ext, ".log%d", &num)
	return num
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs", "app", 1024*10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			i,
			strings.Repeat("X", 50))
		logger.Write([]byte(msg))
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

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compressOld bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.currentFile.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

    if err := os.Rename(r.basePath, rotatedPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    if r.compressOld {
        compressedPath := rotatedPath + ".gz"
        if err := compressFile(rotatedPath, compressedPath); err != nil {
            return fmt.Errorf("failed to compress log file: %w", err)
        }
        os.Remove(rotatedPath)
        rotatedPath = compressedPath
    }

    if err := r.cleanupOldBackups(); err != nil {
        return fmt.Errorf("failed to cleanup old backups: %w", err)
    }

    return r.openCurrentFile()
}

func (r *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return fmt.Errorf("failed to stat log file: %w", err)
    }

    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) cleanupOldBackups() error {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return fmt.Errorf("failed to read directory: %w", err)
    }

    var backupFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            backupFiles = append(backupFiles, filepath.Join(dir, name))
        }
    }

    if len(backupFiles) <= r.maxBackups {
        return nil
    }

    sortBackupFiles(backupFiles)

    for i := 0; i < len(backupFiles)-r.maxBackups; i++ {
        if err := os.Remove(backupFiles[i]); err != nil {
            return fmt.Errorf("failed to remove old backup: %w", err)
        }
    }

    return nil
}

func sortBackupFiles(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractTimestamp(files[i]) > extractTimestamp(files[j]) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(path string) int64 {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    if len(parts) < 2 {
        return 0
    }

    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }

    timestamp, err := strconv.ParseInt(strings.ReplaceAll(timestampStr, "_", ""), 10, 64)
    if err != nil {
        return 0
    }
    return timestamp
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

func (r *LogRotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}