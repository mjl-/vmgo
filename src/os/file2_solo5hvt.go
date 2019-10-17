// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"errors"
	"internal/poll"
	"io"
	"runtime"
	"syscall"
)

var errFS = errors.New("almost no file system")

// fixLongPath is a noop on non-Windows platforms.
func fixLongPath(path string) string {
	return path
}

func rename(oldname, newname string) error {
	return errFS
}

// file is the real representation of *File.
// The extra level of indirection ensures that no clients of os
// can overwrite this data, which could cause the finalizer
// to close the wrong file descriptor.
type file struct {
	pfd         poll.FD
	name        string
	dirinfo     *dirInfo // nil unless directory being read
	nonblock    bool     // whether we set nonblocking mode
	stdoutOrErr bool     // whether this is stdout or stderr
	appendMode  bool     // whether file is opened for appending

	// if set, this is a builtin static file.
	fake *fakeFile
}

type fakeFile struct {
	Path   string
	Offset int64
}

func (f *File) isFake() bool {
	return f.fake != nil
}

// Fd returns the integer Unix file descriptor referencing the open file.
// The file descriptor is valid only until f.Close is called or f is garbage collected.
// On Unix systems this will cause the SetDeadline methods to stop working.
func (f *File) Fd() uintptr {
	if f == nil {
		return ^(uintptr(0))
	}

	if f.isFake() {
		return ^(uintptr(0))
	}

	panic("File.Fd")
}

// NewFile returns a new File with the given file descriptor and
// name. The returned value will be nil if fd is not a valid file
// descriptor. On Unix systems, if the file descriptor is in
// non-blocking mode, NewFile will attempt to return a pollable File
// (one for which the SetDeadline methods work).
func NewFile(fd uintptr, name string) *File {
	switch fd {
	case syscall.Stdin, syscall.Stdout, syscall.Stderr:
	default:
		return nil
	}
	kind := kindNewFile
	return newFile(fd, name, kind)
}

// newFileKind describes the kind of file to newFile.
type newFileKind int

const (
	kindNewFile newFileKind = iota
	kindOpenFile
	kindPipe
	kindNonBlock
)

// newFile is like NewFile, but if called from OpenFile or Pipe
// (as passed in the kind parameter) it tries to add the file to
// the runtime poller.
func newFile(fd uintptr, name string, kind newFileKind) *File {
	fdi := int(fd)
	if fdi < 0 {
		return nil
	}
	f := &File{&file{
		pfd:         poll.FD{},
		name:        name,
		stdoutOrErr: fdi == 1 || fdi == 2,
	}}

	runtime.SetFinalizer(f.file, (*file).close)
	return f
}

// epipecheck raises SIGPIPE if we get an EPIPE error on standard
// output or standard error. See the SIGPIPE docs in os/signal, and
// issue 11845.
func epipecheck(file *File, e error) {
	if e == syscall.EPIPE && file.stdoutOrErr {
		sigpipe()
	}
}

// DevNull is the name of the operating system's ``null device.''
// On Unix-like systems, it is "/dev/null"; on Windows, "NUL".
const DevNull = "/dev/null"

// openFileNolog is the Unix implementation of OpenFile.
// Changes here should be reflected in openFdAt, if relevant.
func openFileNolog(name string, flag int, perm FileMode) (*File, error) {
	if flag != O_RDONLY {
		return nil, errFS
	}

	_, ok := fakeFiles[name]
	if !ok {
		return nil, ErrNotExist
	}

	f := &File{
		file: &file{
			fake: &fakeFile{
				Path:   name,
				Offset: 0,
			},
		},
	}
	return f, nil
}

// Close closes the File, rendering it unusable for I/O.
// On files that support SetDeadline, any pending I/O operations will
// be canceled and return immediately with an error.
// Close will return an error if it has already been called.
func (f *File) Close() error {
	if f == nil {
		return ErrInvalid
	}
	if f.isFake() {
		return nil
	}
	return f.file.close()
}

func (file *file) close() error {
	if file == nil {
		return syscall.EINVAL
	}
	if file.dirinfo != nil {
		file.dirinfo.close()
	}

	// no need for a finalizer anymore
	runtime.SetFinalizer(file, nil)
	return nil
}

// read reads up to len(b) bytes from the File.
// It returns the number of bytes read and an error, if any.
func (f *File) read(b []byte) (n int, err error) {
	if f == nil {
		return 0, syscall.EINVAL
	}

	if f.isFake() {
		n, err = f.pread(b, f.fake.Offset)
		if n > 0 {
			f.fake.Offset += int64(n)
		}
		return
	}

	return 0, syscall.ENOTSUP
}

// pread reads len(b) bytes from the File starting at byte offset off.
// It returns the number of bytes read and the error, if any.
// EOF is signaled by a zero count with err set to nil.
func (f *File) pread(b []byte, off int64) (n int, err error) {
	if f.isFake() {
		buf := fakeFiles[f.fake.Path]
		s := off
		e := s + int64(len(b))
		if e > int64(len(buf)) {
			e = int64(len(buf))
		}
		if s > e {
			s = e
		}
		n = int(e - s)
		copy(b, buf[int(s):int(e)])
		var err error
		if n == 0 {
			err = io.EOF
		}
		return n, err
	}

	return 0, syscall.ENOTSUP
}

// write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
func (f *File) write(b []byte) (n int, err error) {
	if f.isFake() {
		return 0, syscall.ENOSYS
	}

	if f.file.stdoutOrErr {
		runtime.Solo5Write(b)
		return len(b), nil
	}
	return 0, syscall.ENOTSUP
}

// pwrite writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
func (f *File) pwrite(b []byte, off int64) (n int, err error) {
	if f.isFake() {
		return 0, syscall.ENOSYS
	}

	return 0, syscall.ENOTSUP
}

// seek sets the offset for the next Read or Write on file to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error, if any.
func (f *File) seek(offset int64, whence int) (ret int64, err error) {
	if f.isFake() {
		if whence == 0 {
			f.fake.Offset = offset
		} else if whence == 1 {
			f.fake.Offset = f.fake.Offset + offset
		} else {
			buf := fakeFiles[f.fake.Path]
			f.fake.Offset = int64(len(buf)) + offset
		}
		return f.fake.Offset, nil
	}

	return 0, syscall.ENOTSUP
}

// Truncate changes the size of the named file.
// If the file is a symbolic link, it changes the size of the link's target.
// If there is an error, it will be of type *PathError.
func Truncate(name string, size int64) error {
	return syscall.ENOSYS
}

// Remove removes the named file or (empty) directory.
// If there is an error, it will be of type *PathError.
func Remove(name string) error {
	return syscall.ENOSYS
}

func tempDir() string {
	return "/"
}

// Link creates newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.
func Link(oldname, newname string) error {
	return syscall.ENOSYS
}

// Symlink creates newname as a symbolic link to oldname.
// If there is an error, it will be of type *LinkError.
func Symlink(oldname, newname string) error {
	return syscall.ENOSYS
}

func (f *File) readdir(n int) (fi []FileInfo, err error) {
	return nil, syscall.ENOSYS
}

// Readlink returns the destination of the named symbolic link.
// If there is an error, it will be of type *PathError.
func Readlink(name string) (string, error) {
	return "", syscall.ENOSYS
}
