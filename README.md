# filelock

[![Build Status](https://api.travis-ci.org/rosenhouse/filelock.png?branch=master)](http://travis-ci.org/rosenhouse/filelock) [![GoDoc](https://godoc.org/github.com/rosenhouse/filelock?status.png)](https://godoc.org/github.com/rosenhouse/filelock)

Basic inter-process synchronization via POSIX filesystem locking

```bash
go get github.com/rosenhouse/filelock
```

## Usage
- See the [demo code](filelock-demo/main.go)


## Simple example

0. Install the demo binary

  ```bash
  go install github.com/rosenhouse/filelock/filelock-demo
  ```

0. Acquire a lock

  ```bash
  filelock-demo /tmp/some/lock
  ```

0. Open another terminal and attempt to acquire the same lock

  ```bash
  filelock-demo /tmp/some/lock
  ```
