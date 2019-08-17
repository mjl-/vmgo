// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd

package time

import (
	"errors"
	"runtime"
	"sync"
	"syscall"
)

var timezoneZipData = []byte{}

// for testing: whatever interrupts a sleep
func interrupt() {
	syscall.Kill(syscall.Getpid(), syscall.SIGCHLD)
}

type openFile struct {
	offset int
}

var openFiles = struct {
	sync.Mutex
	fds    map[uintptr]*openFile
	nextFD uintptr
}{
	fds:    map[uintptr]*openFile{},
	nextFD: 1,
}

func open(name string) (uintptr, error) {
	if name != runtime.GOROOT()+"/lib/time/zoneinfo.zip" {
		return 0, syscall.ENOENT
	}
	openFiles.Lock()
	defer openFiles.Unlock()
	fd := openFiles.nextFD
	openFiles.nextFD++
	openFiles.fds[fd] = &openFile{0}
	return fd, nil
}

func read(fd uintptr, buf []byte) (int, error) {
	openFiles.Lock()
	defer openFiles.Unlock()
	f, ok := openFiles.fds[fd]
	if !ok {
		return -1, syscall.EBADF
	}
	n := readoff(f, buf, f.offset)
	f.offset += n
	return n, nil
}

func closefd(fd uintptr) {
	openFiles.Lock()
	defer openFiles.Unlock()
	_, ok := openFiles.fds[fd]
	if !ok {
		return
	}
	delete(openFiles.fds, fd)
}

func preadn(fd uintptr, buf []byte, off int) error {
	openFiles.Lock()
	defer openFiles.Unlock()
	f, ok := openFiles.fds[fd]
	if !ok {
		return syscall.EBADF
	}
	if off < 0 {
		off = len(timezoneZipData) + off
		if off < 0 {
			return errors.New("bad seek")
		}
	}
	n := readoff(f, buf, off)
	if n < len(buf) {
		return errors.New("short read")
	}
	return nil
}

func readoff(f *openFile, buf []byte, s int) int {
	e := s + len(buf)
	if e > len(timezoneZipData) {
		e = len(timezoneZipData)
	}
	if s > e {
		s = e
	}
	copy(buf, timezoneZipData[int(s):int(e)])
	return int(e - s)
}

// SetTimezoneDB sets the contents of a timezone zip file that is used to fullfil reads for GOROOT/lib/time/zoneinfo.zip.
func SetTimezoneDB(zipData []byte) {
	timezoneZipData = zipData
}
