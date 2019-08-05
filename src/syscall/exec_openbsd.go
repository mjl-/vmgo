// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

import (
	"sync"
)

type SysProcAttr struct {
}

var ForkLock sync.RWMutex // called by internal/poll, leaving this for now

func CloseOnExec(fd int) { }

func SetNonblock(fd int, nonblocking bool) (err error) {
	flag, err := fcntl(fd, F_GETFL, 0)
	if err != nil {
		return err
	}
	if nonblocking {
		flag |= O_NONBLOCK
	} else {
		flag &^= O_NONBLOCK
	}
	_, err = fcntl(fd, F_SETFL, flag)
	return err
}

// Credential holds user and group identities to be assumed
// by a child process started by StartProcess.
type Credential struct {
}

// ProcAttr holds attributes that will be applied to a new process started
// by StartProcess.
type ProcAttr struct {
}


// Combination of fork and exec, careful to be thread safe.
func ForkExec(argv0 string, argv []string, attr *ProcAttr) (pid int, err error) {
	return -1, ENOSYS
}

// StartProcess wraps ForkExec for package os.
func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error) {
	return -1, 0, ENOSYS
}

// Exec invokes the execve(2) system call.
func Exec(argv0 string, argv []string, envv []string) (err error) {
	return ENOSYS
}
