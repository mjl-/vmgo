// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	_ "unsafe"
)

type mOS struct {}
type sigset uint32
type gsignalStack struct{}

func osinit() {
	ncpu = 1
	physPageSize = 4*1024
}

//go:nosplit
func getRandomData(r []byte) {
	// xxx todo bug, need to actually set random data!
}

const _NSIG = 0

func signame(sig uint32) string {
	return ""
}

func initsig(preinit bool) {
}

func goenvs() {
	solo5envs()
}

// Called to initialize a new m (including the bootstrap m).
// Called on the parent thread (main thread in case of bootstrap), can allocate memory.
func mpreinit(mp *m) {
	mp.gsignal = malg(2 * 1024)
	mp.gsignal.m = mp
}

//go:nosplit
func msigsave(mp *m) {
}

//go:nosplit
func msigrestore(sigmask sigset) {
}

//go:nosplit
//go:nowritebarrierrec
func clearSignalHandlers() {
}

//go:nosplit
func sigblock() {
}

// Called to initialize a new m (including the bootstrap m).
// Called on the new thread, can not allocate memory.
func minit() {
}

// Called from dropm to undo the effect of an minit.
//go:nosplit
func unminit() {
}

func crash() {
	*(*int32)(nil) = 0
}

func setProcessCPUProfiler(hz int32) {}
func setThreadCPUProfiler(hz int32)  {}
func sigdisable(uint32)              {}
func sigenable(uint32)               {}
func sigignore(uint32)               {}

func nanotime() int64 {
	return 1
	// xxx
}

func walltime() (sec int64, nsec int32) {
	return 0, 0
	// xxx
}

//go:linkname os_sigpipe os.sigpipe
func os_sigpipe() {
	throw("too many writes on closed pipe")
}
