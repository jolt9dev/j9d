package cps

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/jolt9dev/j9d/pkg/platform"
)

const (
	ARCH     = runtime.GOARCH
	PLATFORM = runtime.GOOS
)

var (
	Args     = os.Args[1:]
	Stderr   = os.Stderr
	Stdout   = os.Stdout
	Stdin    = os.Stdin
	ExecPath = os.Args[0]
	history  = []string{}
	reader   *bufio.Reader
	writer   *bufio.Writer
	eol      = []byte(platform.EOL)
)

func init() {
	current, _ := Cwd()
	history = append(history, current)
}

func Cwd() (string, error) {
	return os.Getwd()
}

func Exit(code int) {
	os.Exit(code)
}

func Pid() int {
	return os.Getpid()
}

func Ppid() int {
	ppid := os.Getppid()
	if ppid == 0 {
		return -1
	}

	return ppid
}

func Uid() int {
	return os.Getuid()
}

func Gid() int {
	return os.Getgid()
}

func Euid() int {
	return os.Geteuid()
}

func Egid() int {
	return os.Getegid()
}

func Pushd(path string) error {
	history = append(history, path)
	return os.Chdir(path)
}

func Popd() error {
	if len(history) == 1 {
		return nil
	}

	last := history[len(history)-1]
	history = history[:len(history)-1]
	return os.Chdir(last)
}

func Read(b []byte) (int, error) {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	return reader.Read(b)
}

func ReadLine() (string, error) {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	b, _, e := reader.ReadLine()
	return string(b), e
}

func WriteBytes(b []byte) (int, error) {
	return Stdout.Write(b)
}

func WriteRune(r rune) (int, error) {
	if writer == nil {
		writer = bufio.NewWriter(os.Stdout)
	}

	return writer.WriteRune(r)
}

func WriteString(s string) (int, error) {
	if writer == nil {
		writer = bufio.NewWriter(Stdout)
	}

	b, err := writer.WriteString(s)
	if err != nil {
		return b, err
	}
	writer.Flush()
	return b, err
}

func Writef(format string, a ...interface{}) (int, error) {
	msg := fmt.Sprintf(format, a...)
	return WriteString(msg)
}

func Writeln(s string) (int, error) {

	n, err := WriteString(s)
	if err != nil {
		return n, err
	}

	n2, err := writer.Write(eol)
	return n + n2, err
}
