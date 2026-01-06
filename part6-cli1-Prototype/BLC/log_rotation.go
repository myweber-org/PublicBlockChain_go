
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