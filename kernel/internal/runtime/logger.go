package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	infoPath  string
	errorPath string
	mu        sync.Mutex
}

func NewLogger(paths Paths) *Logger {
	return &Logger{
		infoPath:  filepath.Join(paths.CurrentLogDir, "info.log"),
		errorPath: filepath.Join(paths.CurrentLogDir, "error.log"),
	}
}

func (l *Logger) Info(msg string) {
	l.write("INFO", msg, "")
}

func (l *Logger) Error(msg string, err error) {
	detail := ""
	if err != nil {
		detail = err.Error()
	}
	l.write("ERROR", msg, detail)
}

func (l *Logger) write(level, msg, detail string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := fmt.Sprintf("[%s] [%s] %s", time.Now().Format(time.RFC3339), level, msg)
	if detail != "" {
		entry += "\n  Detail: " + detail
	}
	entry += "\n"

	_ = os.WriteFile(l.infoPath, []byte{}, FilePerm())
	_ = appendFile(l.infoPath, entry)
	if level == "ERROR" {
		_ = os.WriteFile(l.errorPath, []byte{}, FilePerm())
		_ = appendFile(l.errorPath, entry)
	}
}

func appendFile(path, entry string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FilePerm())
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(entry)
	return err
}
