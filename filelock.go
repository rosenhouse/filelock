package filelock

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
)

type Locker struct {
	Path string
}

type LockedFile interface {
	io.ReadWriteSeeker
	io.Closer
}

func (l *Locker) Open() (LockedFile, error) {
	dir := filepath.Dir(l.Path)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		panic(err)
	}
	const flags = os.O_RDWR | os.O_CREATE | syscall.LOCK_EX
	return os.OpenFile(l.Path, flags, 0600)
}
