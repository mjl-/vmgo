// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd

package time

import (
	"syscall"
)

// for testing: whatever interrupts a sleep
func interrupt() {
	syscall.Kill(syscall.Getpid(), syscall.SIGCHLD)
}

func open(name string) (uintptr, error) {
	return 0, syscall.ENOTSUP
}

func read(fd uintptr, buf []byte) (int, error) {
	return -1, syscall.ENOTSUP
}

func closefd(fd uintptr) {
}

func preadn(fd uintptr, buf []byte, off int) error {
	return syscall.ENOTSUP
}
