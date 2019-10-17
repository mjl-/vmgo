// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
	"time"
)

// syscallMode returns the syscall-specific mode bits from Go's portable mode bits.
func syscallMode(i FileMode) (o uint32) {
	o |= uint32(i.Perm())
	if i&ModeSetuid != 0 {
		o |= syscall.S_ISUID
	}
	if i&ModeSetgid != 0 {
		o |= syscall.S_ISGID
	}
	if i&ModeSticky != 0 {
		o |= syscall.S_ISVTX
	}
	// No mapping for Go's ModeTemporary (plan9 only).
	return
}

func sigpipe() // implemented in package runtime

// See docs in file.go:Chmod.
func chmod(name string, mode FileMode) error {
	return syscall.ENOSYS
}

// See docs in file.go:(*File).Chmod.
func (f *File) chmod(mode FileMode) error {
	return syscall.ENOSYS
}

// Chown changes the numeric uid and gid of the named file.
// If the file is a symbolic link, it changes the uid and gid of the link's target.
// A uid or gid of -1 means to not change that value.
// If there is an error, it will be of type *PathError.
//
// On Windows or Plan 9, Chown always returns the syscall.EWINDOWS or
// EPLAN9 error, wrapped in *PathError.
func Chown(name string, uid, gid int) error {
	return syscall.ENOSYS
}

// Lchown changes the numeric uid and gid of the named file.
// If the file is a symbolic link, it changes the uid and gid of the link itself.
// If there is an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped
// in *PathError.
func Lchown(name string, uid, gid int) error {
	return syscall.ENOSYS
}

// Chown changes the numeric uid and gid of the named file.
// If there is an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped
// in *PathError.
func (f *File) Chown(uid, gid int) error {
	return syscall.ENOSYS
}

// Truncate changes the size of the file.
// It does not change the I/O offset.
// If there is an error, it will be of type *PathError.
func (f *File) Truncate(size int64) error {
	return syscall.ENOSYS
}

// Sync commits the current contents of the file to stable storage.
// Typically, this means flushing the file system's in-memory copy
// of recently written data to disk.
func (f *File) Sync() error {
	return syscall.ENOSYS
}

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
// If there is an error, it will be of type *PathError.
func Chtimes(name string, atime time.Time, mtime time.Time) error {
	return syscall.ENOSYS
}

// Chdir changes the current working directory to the file,
// which must be a directory.
// If there is an error, it will be of type *PathError.
func (f *File) Chdir() error {
	return syscall.ENOSYS
}

// setDeadline sets the read and write deadline.
func (f *File) setDeadline(t time.Time) error {
	if f.isFake() {
		return syscall.ENOSYS
	}

	panic("File.setDeadline")
	return syscall.ENOSYS
}

// setReadDeadline sets the read deadline.
func (f *File) setReadDeadline(t time.Time) error {
	if f.isFake() {
		return syscall.ENOSYS
	}

	panic("File.setReadDeadline")
	return syscall.ENOSYS
}

// setWriteDeadline sets the write deadline.
func (f *File) setWriteDeadline(t time.Time) error {
	if f.isFake() {
		return syscall.ENOSYS
	}

	panic("File.setWriteDeadline")
	return syscall.ENOSYS
}

// checkValid checks whether f is valid for use.
// If not, it returns an appropriate error, perhaps incorporating the operation name op.
func (f *File) checkValid(op string) error {
	if f == nil {
		return ErrInvalid
	}
	return nil
}

type rawConn struct{}

func (c *rawConn) Control(f func(uintptr)) error {
	return syscall.ENOSYS
}

func (c *rawConn) Read(f func(uintptr) bool) error {
	return syscall.ENOSYS
}

func (c *rawConn) Write(f func(uintptr) bool) error {
	return syscall.ENOSYS
}

func newRawConn(file *File) (*rawConn, error) {
	return nil, syscall.ENOSYS
}
