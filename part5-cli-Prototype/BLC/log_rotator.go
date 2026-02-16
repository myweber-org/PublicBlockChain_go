package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize   = 1024 * 1024 // 1MB
	maxArchives  = 5
	logFileName  = "app.log"
	archiveDir   = "archives"
)

func rotateLog() error {
	info, err := os.Stat(logFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() < maxLogSize {
		return nil
	}

	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveName := filepath.Join(archiveDir, fmt.Sprintf("%s_%s.log", logFileName, timestamp))

	if err := os.Rename(logFileName, archiveName); err != nil {
		return err
	}

	cleanupOldArchives()
	return nil
}

func cleanupOldArchives() {
	files, err := filepath.Glob(filepath.Join(archiveDir, "*.log"))
	if err != nil {
		return
	}

	if len(files) <= maxArchives {
		return
	}

	for i := 0; i < len(files)-maxArchives; i++ {
		os.Remove(files[i])
	}
}

func writeLog(message string) error {
	if err := rotateLog(); err != nil {
		return err
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	_, err = io.WriteString(file, logEntry)
	return err
}

func main() {
	for i := 1; i <= 100; i++ {
		message := fmt.Sprintf("Log entry number %d", i)
		if err := writeLog(message); err != nil {
			fmt.Printf("Error writing log: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("Log rotation demonstration completed")
}
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentSize int64
    file        *os.File
    basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    fullPath := filepath.Join(path, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    stat, err := file.Stat()
    if err != nil {
        return nil, err
    }

    return &RotatingLogger{
        currentSize: stat.Size(),
        file:        file,
        basePath:    path,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s.gz", logFileName, timestamp)
    backupPath := filepath.Join(rl.basePath, backupName)

    if err := compressFile(filepath.Join(rl.basePath, logFileName), backupPath); err != nil {
        return err
    }

    if err := cleanupOldBackups(rl.basePath); err != nil {
        log.Printf("Failed to cleanup old backups: %v", err)
    }

    file, err := os.OpenFile(filepath.Join(rl.basePath, logFileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    return nil
}

func compressFile(source, target string) error {
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

func cleanupOldBackups(dir string) error {
    pattern := filepath.Join(dir, logFileName+".*.gz")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackupFiles {
        return nil
    }

    for i := 0; i < len(matches)-maxBackupFiles; i++ {
        if err := os.Remove(matches[i]); err != nil {
            return err
        }
    }
    return nil
}

func (rl *RotatingLogger) Close() error {
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(logger)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(10 * time.Millisecond)
    }
}