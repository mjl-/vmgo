// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
)

// Auxiliary information if the File describes a directory
type dirInfo struct {
	buf  []byte // buffer for directory I/O
	nbuf int    // length of buf; return value from Getdirentries
	bufp int    // location of next record in buf.
}

const (
	// More than 5760 to work around https://golang.org/issue/24015.
	blockSize = 8192
)

func (d *dirInfo) close() {}

func (f *File) readdirnames(n int) (names []string, err error) {
	return nil, syscall.ENOSYS
}
