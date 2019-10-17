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

// The only signal values guaranteed to be present in the os package on all
// systems are os.Interrupt (send the process an interrupt) and os.Kill (force
// the process to exit). On Windows, sending os.Interrupt to a process with
// os.Process.Signal is not implemented; it will return an error instead of
// sending a signal.
var (
	Interrupt Signal = syscall.SIGINT
	Kill      Signal = syscall.SIGKILL
)

func startProcess(name string, argv []string, attr *ProcAttr) (p *Process, err error) {
	return nil, syscall.ENOTSUP
}

func (p *Process) kill() error {
	return syscall.ENOTSUP
}

// ProcessState stores information about a process, as reported by Wait.
type ProcessState struct {
}

// Pid returns the process id of the exited process.
func (p *ProcessState) Pid() int {
	return 0
}

func (p *ProcessState) exited() bool {
	return false
}

func (p *ProcessState) success() bool {
	return false
}

func (p *ProcessState) sys() interface{} {
	return nil
}

func (p *ProcessState) sysUsage() interface{} {
	return nil
}

func (p *ProcessState) String() string {
	return "<nil>"
}

// ExitCode returns the exit code of the exited process, or -1
// if the process hasn't exited or was terminated by a signal.
func (p *ProcessState) ExitCode() int {
	return -1
}
