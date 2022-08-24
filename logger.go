package golog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	// LevelDebug log level
	LevelDebug = 1 << iota
	// LevelInfo log level
	LevelInfo
	// LevelWarning
	LevelWarning
	// LevelError log level
	LevelError
	// LevelCritical log level
	LevelCritical
)

var (
	DefaultLayout = time.RFC3339Nano
	DefaultLevel  = LevelInfo
)

type Logger struct {
	sync.Mutex
	w     io.Writer
	mock  bool
	close func() error

	Level        int
	Layout       string
	ErrorHandler func(error)
}

func OpenLogFile(path string, mode os.FileMode) (fd *os.File, err error) {

	if fd, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, mode); err != nil {
		return
	}

	if _, err = fd.Seek(0, os.SEEK_END); err != nil {
		return
	}

	return
}

func FromWriter(w io.Writer) *Logger {
	return &Logger{
		w:      w,
		Level:  LevelInfo,
		Layout: DefaultLayout}
}

func FromWriteCloser(w io.WriteCloser) *Logger {
	l := FromWriter(w)
	l.close = w.Close

	return l
}

func FromStdout() *Logger {
	return FromWriter(os.Stdout)
}

func FromFile(fd *os.File) (l *Logger) {
	return FromWriteCloser(fd)
}

func FromPath(path string, mode os.FileMode) (l *Logger, err error) {
	var lf *os.File

	if lf, err = OpenLogFile(path, mode); err != nil {
		return
	}

	l = FromWriteCloser(lf)
	return
}

func (l *Logger) writeString(s string) (err error) {
	_, err = l.w.Write([]byte(s))
	return
}

func (l *Logger) output(prefix string, i ...interface{}) string {
	return fmt.Sprintf(l.format(prefix, i...), i...)
}

func (l *Logger) log(msg string) (err error) {
	l.Lock()
	defer l.Unlock()
	return l.writeString(msg)
}

func (l *Logger) timestamp() []byte {
	return append(append([]byte("["), []byte(time.Now().Format(l.Layout))...), ']')
}

func (l *Logger) format(prefix string, i ...interface{}) string {
	fmt := make([][]byte, len(i)+2)
	fmt[0] = l.timestamp()
	fmt[1] = []byte(prefix)
	for i := 2; i < len(fmt); i++ {
		fmt[i] = []byte("%v")
	}
	return string(append(bytes.Join(fmt, []byte(" ")), '\n'))
}

func (l *Logger) handleError(err error) {
	if l.ErrorHandler != nil {
		l.ErrorHandler(err)
	}
}

func (l *Logger) Debug(i ...interface{}) {
	if l.Level <= LevelDebug {
		msg := l.output("DEBUG -", i...)
		l.log(msg)
	}
}

func (l *Logger) Debugf(format string, i ...interface{}) {
	l.Debug(fmt.Sprintf(format, i...))
}

func (l *Logger) Info(i ...interface{}) {
	if l.Level <= LevelInfo {
		msg := l.output("INFO -", i...)
		l.log(msg)
	}
}

func (l *Logger) Infof(format string, i ...interface{}) {
	l.Info(fmt.Sprintf(format, i...))
}

func (l *Logger) Warn(i ...interface{}) {
	if l.Level <= LevelWarning {
		msg := l.output("WARNING -", i...)
		l.log(msg)
	}
}

func (l *Logger) Warnf(format string, i ...interface{}) {
	l.Warn(fmt.Sprintf(format, i...))
}

func (l *Logger) Error(i ...interface{}) {
	if l.Level <= LevelError {
		msg := l.output("ERROR -", i...)
		l.log(msg)
		l.handleError(errors.New(msg))
	}
}

func (l *Logger) Errorf(format string, i ...interface{}) {
	l.Error(fmt.Sprintf(format, i...))
}

func (l *Logger) Critical(i ...interface{}) {
	if l.Level <= LevelCritical {
		msg := l.output("CRITICAL -", i...)
		l.log(msg)
		l.handleError(errors.New(msg))
	}
}

func (l *Logger) Criticalf(format string, i ...interface{}) {
	l.Critical(fmt.Sprintf(format, i...))
}

func (l *Logger) Abort(rc int, i ...interface{}) {
	l.log(l.output("ABORT -", i...))
	if !l.mock {
		os.Exit(rc)
	}
}

func (l *Logger) Close() error {
	if l.close != nil {
		l.Lock()
		defer l.Unlock()
		return l.close()
	}
	return nil
}
