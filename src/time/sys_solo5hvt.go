// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"errors"
	"sync"
	"syscall"
)

var timezoneZipData = []byte{}

type openFile struct {
	offset int
}

var openFiles = struct {
	sync.Mutex
	fds				map[uintptr]*openFile
	nextFD uintptr
}{
	fds:				map[uintptr]*openFile{},
	nextFD: 1,
}

func open(name string) (uintptr, error) {
	if name != "/zoneinfo.zip" || len(timezoneZipData) == 0 {
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

// SetTimezoneDB sets the contents of the timezone zip file used for looking up time zones.
func SetTimezoneDB(zipData []byte) {
	timezoneZipData = zipData
}
