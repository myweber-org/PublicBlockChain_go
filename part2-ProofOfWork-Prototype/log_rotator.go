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

type Rotator struct {
	filePath string
	current  *os.File
	size     int64
}

func NewRotator(path string) (*Rotator, error) {
	r := &Rotator{filePath: path}
	if err := r.openFile(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
	if r.size+int64(len(p)) > maxFileSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := r.current.Write(p)
	if err == nil {
		r.size += int64(n)
	}
	return n, err
}

func (r *Rotator) rotate() error {
	if r.current != nil {
		r.current.Close()
		backupPath := fmt.Sprintf("%s.%s.gz", r.filePath, time.Now().Format("20060102150405"))
		if err := compressFile(r.filePath, backupPath); err != nil {
			return err
		}
		os.Remove(r.filePath)
		cleanOldBackups(r.filePath)
	}
	return r.openFile()
}

func (r *Rotator) openFile() error {
	f, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	r.current = f
	r.size = info.Size()
	return nil
}

func compressFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
}

func cleanOldBackups(basePath string) {
	pattern := filepath.Join(filepath.Dir(basePath), filepath.Base(basePath)+".*.gz")
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

func (r *Rotator) Close() error {
	if r.current != nil {
		return r.current.Close()
	}
	return nil
}