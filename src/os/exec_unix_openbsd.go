// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
	"time"
)

func (p *Process) wait() (ps *ProcessState, err error) {
	return nil, syscall.ENOTSUP
}

func (p *Process) signal(sig Signal) error {
	return syscall.ENOTSUP
}

func (p *Process) release() error {
	return syscall.ENOTSUP
}

func findProcess(pid int) (p *Process, err error) {
	return nil, syscall.ENOTSUP
}

func (p *ProcessState) userTime() time.Duration {
	return 0
}

func (p *ProcessState) systemTime() time.Duration {
	return 0
}
