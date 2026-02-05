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