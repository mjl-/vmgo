// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
	"time"
)

var zeroTime time.Time

// Stat returns the FileInfo structure describing file.
// If there is an error, it will be of type *PathError.
func (f *File) Stat() (FileInfo, error) {
	if f == nil {
		return nil, ErrInvalid
	}
	if f.isFake() {
		return fakeStat(f.fake.Path)
	}
	return nil, syscall.ENOSYS
}

func fakeStat(path string) (FileInfo, error) {
	buf, ok := fakeFiles[path]
	if !ok {
		return nil, ErrNotExist
	}
	st := &fileStat{
		name: path,
		size: int64(len(buf)),
		mode: 0,
		modTime: zeroTime,
		sys: syscall.Stat_t{},
	}
	return st, nil
}

// statNolog stats a file with no test logging.
func statNolog(name string) (FileInfo, error) {
	return fakeStat(name)
}

// lstatNolog lstats a file with no test logging.
func lstatNolog(name string) (FileInfo, error) {
	return fakeStat(name)
}
