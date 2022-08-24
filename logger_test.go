package golog

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/0xrawsec/toast"
)

func TestLogger(t *testing.T) {
	tt := toast.FromT(t)

	l := FromStdout()
	defer l.Close()
	l.mock = true
	l.ErrorHandler = func(err error) { panic(err) }

	l.Debug("we are debugging this buggy code")
	l.Debugf("why this code doesn't work ? %d == %d must be %t", 42, 42, true)
	l.Info("very useful information")
	l.Infof("%s %d > %d", "do you know that", 42, 43)
	tt.ShouldPanic(func() { l.Error("noooooooooooo") })
	tt.ShouldPanic(func() { l.Errorf("%s %d %t", "ahhhhhhh", 2, true) })
	tt.ShouldPanic(func() { l.Critical("noooooooooooo") })
	tt.ShouldPanic(func() { l.Criticalf("%s %d %t", "ahhhhhhh", 2, true) })
	l.Abort(42, "aborting")
	l.Logf("simple log without prefix")
}

func bufferLogger() (buf *bytes.Buffer, l *Logger) {
	buf = new(bytes.Buffer)
	l = FromWriter(buf)
	l.mock = true
	return
}

func TestDebug(t *testing.T) {
	tt := toast.FromT(t)

	b, l := bufferLogger()
	l.Level = LevelDebug
	l.Debugf("dummy message")
	tt.Assert(strings.Contains(b.String(), " DEBUG "))
	b.Reset()
	l.Infof("dummy message")
	tt.Assert(strings.Contains(b.String(), " INFO "))
	b.Reset()
	l.Warnf("dummy message")
	tt.Assert(strings.Contains(b.String(), " WARNING "))
	b.Reset()
	l.Criticalf("dummy message")
	tt.Assert(strings.Contains(b.String(), " CRITICAL "))
	b.Reset()
	l.Abort(42, "dummy message")
	tt.Assert(strings.Contains(b.String(), " ABORT "))
	b.Reset()
	l.Logf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " - "))
}

func TestInfo(t *testing.T) {
	tt := toast.FromT(t)

	b, l := bufferLogger()
	l.Level = LevelInfo
	l.Debugf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " DEBUG "))
	b.Reset()
	l.Infof("dummy message")
	tt.Assert(strings.Contains(b.String(), " INFO "))
	b.Reset()
	l.Warnf("dummy message")
	tt.Assert(strings.Contains(b.String(), " WARNING "))
	b.Reset()
	l.Criticalf("dummy message")
	tt.Assert(strings.Contains(b.String(), " CRITICAL "))
	b.Reset()
	l.Abort(42, "dummy message")
	tt.Assert(strings.Contains(b.String(), " ABORT "))
	b.Reset()
	l.Logf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " - "))
}

func TestWarning(t *testing.T) {
	tt := toast.FromT(t)

	b, l := bufferLogger()
	l.Level = LevelWarning
	l.Debugf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " DEBUG "))
	b.Reset()
	l.Infof("dummy message")
	tt.Assert(!strings.Contains(b.String(), " INFO "))
	b.Reset()
	l.Warnf("dummy message")
	tt.Assert(strings.Contains(b.String(), " WARNING "))
	b.Reset()
	l.Criticalf("dummy message")
	tt.Assert(strings.Contains(b.String(), " CRITICAL "))
	b.Reset()
	l.Abort(42, "dummy message")
	tt.Assert(strings.Contains(b.String(), " ABORT "))
	b.Reset()
	l.Logf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " - "))
}

func TestCritical(t *testing.T) {
	tt := toast.FromT(t)

	b, l := bufferLogger()
	l.Level = LevelCritical
	l.Debugf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " DEBUG "))
	b.Reset()
	l.Infof("dummy message")
	tt.Assert(!strings.Contains(b.String(), " INFO "))
	b.Reset()
	l.Warnf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " WARNING "))
	b.Reset()
	l.Criticalf("dummy message")
	tt.Assert(strings.Contains(b.String(), " CRITICAL "))
	b.Reset()
	l.Abort(42, "dummy message")
	tt.Assert(strings.Contains(b.String(), " ABORT "))
	b.Reset()
	l.Logf("dummy message")
	tt.Assert(!strings.Contains(b.String(), " - "))
}

func readlines(path string) (lines []string) {
	fd, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fileScanner := bufio.NewScanner(fd)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return
}

func TestOpenPath(t *testing.T) {
	tt := toast.FromT(t)

	path := filepath.Join(t.TempDir(), "logfile.log")
	l, err := FromPath(path, 0777)
	tt.CheckErr(err)

	l.Info("dummy")
	tt.CheckErr(l.Close())

	l, err = FromPath(path, 0777)
	tt.CheckErr(err)
	defer l.Close()

	l.Warn("dummy")
	l.Error("dummy")

	lines := readlines(path)
	tt.Assert(strings.Contains(lines[0], " INFO "))
	tt.Assert(strings.Contains(lines[1], " WARNING "))
	tt.Assert(strings.Contains(lines[2], " ERROR "))
}
