// Credits to https://stackoverflow.com/a/28797984
// Interface to io.Writer to automaticly rotate the
// log file based on max size and/or time
package auralog

import (
  "os"
  "sync"
  "time"
  "strings"
  "path/filepath"
)

var (
  Megabyte int64 = 1024 * 1024 // Variable to use along with MaxSize to use megabytes.
  Kilobyte int64 = 1024 // Variable to use along with MaxSize to use kilobytes.
)

// The structure for RotateWriter, which should interface io.Writer
type RotateWriter struct {
  Dir string // the directory to put log files.
  Filename string // should be set to the actual filename and extension.
  ExTime time.Duration // how often the log should rotate.
  MaxSize int64 // max size a log file is allowed to be in bytes.

  lock sync.Mutex
  now time.Time
  fp *os.File
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
  w.lock.Lock()
  defer w.lock.Unlock()

  if w.fp == nil {
    if err := w.Resume(); err != nil {
      return 0, err
    }
  }

  if w.now.After(w.now.Add(w.ExTime)) {
    if err := w.Rotate(); err != nil {
      return 0, err
    }
  }

  return w.fp.Write(output)
}

// Resume the current log file instead of creating a new one each time program is started back up.
func (w *RotateWriter) Resume() error {
  var err error
  var filename = w.Dir+w.Filename

  // Check if file already exixts to resume it.
  if _, err := os.Stat(filename); err != nil {
    return err
  }

  w.now = time.Now()
  w.fp, err = os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0666)
  return err
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) Rotate() error {
  var err error
  var filename = w.Dir+w.Filename

  // create the needed direactories if they don't exists.
  if err := os.MkdirAll(w.Dir, 0755); err != nil {
    return err
  }

  // Close existing file if open
  if w.fp != nil {
    if err := w.fp.Close(); err != nil {
      return err
    }
    w.fp = nil
  }

  // Rename dest file if it already exists
  if fi, err := os.Stat(filename); err == nil {
    if err := w.renameFile(); err != nil {
      return err
    }
  }else{
    if w.MaxSize > 0 {
      if fi.Size() >= w.MaxSize {
        if err := w.renameFile(); err != nil {
          return err
        }
      }
    }
  }

  // Create a file.
  w.now = time.Now()
  w.fp, err = os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0666)
  return err
}

// Rename the log file to include the current date. Uses RFC3339 time format.
func (w *RotateWriter) renameFile() error {
  var filename = w.Dir+w.Filename
  newfn := filename[:len(filename)-len(filepath.Ext(w.Filename))]+"-"+time.Now().Format(time.RFC3339)+filepath.Ext(w.Filename)
  err := os.Rename(filename, cleanName(newfn));
  return err
}

// Replace : to _ for os.Rename.
func cleanName(name string) string {
  return strings.ReplaceAll(name, ":", "_")
}