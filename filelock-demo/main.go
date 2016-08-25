package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/rosenhouse/filelock"
)

func main() {
	if err := mainWithError(); err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: %s\n", err)
		os.Exit(1)
	}
}

func mainWithError() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s <file_to_lock>", os.Args[0])
	}

	locker := filelock.Locker{Path: os.Args[1]}

	startTime := time.Now()
	fmt.Fprintf(os.Stderr, "waiting to acquire lock on %s...", locker.Path)

	file, err := locker.Open()
	if err != nil {
		return fmt.Errorf("acquire: %s", err)
	}

	durationToAcquire := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "done after %f milliseconds.\n", durationToAcquire.Seconds())

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read: %s", err)
	}
	fmt.Fprintf(os.Stderr, "printing file contents to stdout:\n")

	os.Stdout.Write(contents)

	fmt.Fprintf(os.Stderr, "\nEnter new file contents and then Ctrl+D (EOF) to release the lock\n")

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seek: %s", err)
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncate: %s", err)
	}

	_, err = io.Copy(file, os.Stdin)
	if err != nil {
		return fmt.Errorf("read stdin: %s", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("close: %s", err)
	}

	fmt.Fprintf(os.Stderr, "Complete!\n")

	return nil
}