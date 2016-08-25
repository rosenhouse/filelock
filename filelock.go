package filelock

import (
	"os"
	"path/filepath"
	"syscall"
)

type Locker struct {
	Path string
}

func (l *Locker) Open() (*os.File, error) {
	dir := filepath.Dir(l.Path)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		panic(err)
	}
	const flags = os.O_RDWR | os.O_CREATE | syscall.LOCK_EX
	file, err := os.OpenFile(l.Path, flags, 0600)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		return nil, err
	}
	return file, nil
}
